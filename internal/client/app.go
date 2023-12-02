package client

import (
	"context"
	"encoding/json"
	"fmt"
	"k8s-node-latency/internal/server"
	"log"
	"net/http"
	"net/http/httptrace"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	serviceName = "node-latency-server-service"
	namespace   = "node-latency"
)

const (
	EnvKubePodName  = "KUBE_POD_NAME"
	EnvKubeNodeName = "KUBE_NODE_NAME"
	RoundDelay      = 500 * time.Millisecond
)

type Client struct {
	podName       string
	nodeName      string
	serverAddress string

	kubeClient *kubernetes.Clientset
}

func NewClient() *Client {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Panic("couldn't fetch the in-cluster config", err)
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panic("couldn't create kubernetes config", err)
	}

	service, err := kubeClient.CoreV1().Services(namespace).Get(context.Background(), serviceName, metav1.GetOptions{})
	if err != nil {
		log.Panic("couldn't get server Service: ", err)
	}

	clusterIP := service.Spec.ClusterIP
	log.Printf("ClusterIP of service %s in namespace %s is %s\n", serviceName, namespace, clusterIP)

	podName, exists := os.LookupEnv(EnvKubePodName)
	if !exists {
		log.Panic("couldn't retrieve podname")
	}

	nodeName, exists := os.LookupEnv(EnvKubeNodeName)
	if !exists {
		log.Panic("couldn't retrieve node name")
	}

	return &Client{
		nodeName:      nodeName,
		podName:       podName,
		serverAddress: clusterIP,
		kubeClient:    kubeClient,
	}
}

func (c *Client) testLoop() {
	for {
		c.sendHTTPPing()
		time.Sleep(RoundDelay)
	}
}

func (c *Client) Run() {
	go c.startMetricsServer()
	go c.testLoop()

	sigint := make(chan os.Signal, 1)
	// interrupt signal sent from terminal
	signal.Notify(sigint, os.Interrupt)
	// sigterm signal sent from kubernetes
	signal.Notify(sigint, syscall.SIGTERM)

	<-sigint
	log.Fatal("SIGINT received, shutting down client")
}

func (c *Client) sendHTTPPing() {
	url := fmt.Sprintf("http://%s:8080/ping", c.serverAddress)
	req, _ := http.NewRequest("GET", url, nil)

	var start, connect, dns time.Time
	var dnsLatency, connLatency, firstByteLatency float64

	trace := &httptrace.ClientTrace{
		DNSStart: func(dsi httptrace.DNSStartInfo) { dns = time.Now() },
		DNSDone: func(ddi httptrace.DNSDoneInfo) {
			dnsLatency = float64(time.Since(dns).Milliseconds())
		},

		ConnectStart: func(network, addr string) { connect = time.Now() },
		ConnectDone: func(network, addr string, err error) {
			connLatency = float64(time.Since(connect).Milliseconds())
		},

		GotFirstResponseByte: func() {
			firstByteLatency = float64(time.Since(start).Milliseconds())
		},
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	start = time.Now()

	// Create a custom transport with DisableKeepAlives set to true
	transport := &http.Transport{
		DisableKeepAlives: true,
	}

	// Use the custom transport in the client
	client := &http.Client{
		Transport: transport,
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	end := time.Since(start)

	defer res.Body.Close()

	var resNodeInfo server.NodeLatencyResponse
	err = json.NewDecoder(res.Body).Decode(&resNodeInfo)
	if err != nil {
		panic(err)
	}

	//log.Printf("One way latency: %f", float64(resNodeInfo.OneTripTime.Sub(start).Milliseconds()))
	//log.Printf("Total latency: %f", float64(end.Milliseconds()))
	log.Printf("Duration of request time: %s", end)
	log.Printf("Duration of one-way trip: %s", resNodeInfo.OneTripTime.Sub(start))
	log.Printf("trip completed: %s -> %s", c.podName, resNodeInfo.ServerPodName)
	dnsLatencyHistogram.
		WithLabelValues(c.podName, c.nodeName, resNodeInfo.ServerPodName, resNodeInfo.ServerNodeName).
		Observe(dnsLatency)
	connectionLatencyHistogram.
		WithLabelValues(c.podName, c.nodeName, resNodeInfo.ServerPodName, resNodeInfo.ServerNodeName).
		Observe(connLatency)
	firstByteLatencyHistogram.
		WithLabelValues(c.podName, c.nodeName, resNodeInfo.ServerPodName, resNodeInfo.ServerNodeName).
		Observe(firstByteLatency)
	httpTotalLatencyHistogram.
		WithLabelValues(c.podName, c.nodeName, resNodeInfo.ServerPodName, resNodeInfo.ServerNodeName).
		Observe(float64(end.Milliseconds()))
	httpTotalLatencySummary.
		WithLabelValues(c.podName, c.nodeName, resNodeInfo.ServerPodName, resNodeInfo.ServerNodeName).
		Observe(float64(end.Milliseconds()))
	httpOneWayTripLatencySummary.
		WithLabelValues(c.podName, c.nodeName, resNodeInfo.ServerPodName, resNodeInfo.ServerNodeName).
		Observe(float64(resNodeInfo.OneTripTime.Sub(start).Milliseconds()))

	// Close the transport if it's no longer needed
	transport.CloseIdleConnections()
}

func (c *Client) startMetricsServer() {
	log.Println("starting client metrics server")

	prometheus.MustRegister(dnsLatencyHistogram)
	prometheus.MustRegister(connectionLatencyHistogram)
	prometheus.MustRegister(firstByteLatencyHistogram)
	prometheus.MustRegister(httpTotalLatencyHistogram)
	prometheus.MustRegister(httpTotalLatencySummary)
	prometheus.MustRegister(httpOneWayTripLatencySummary)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8081", nil))
}
