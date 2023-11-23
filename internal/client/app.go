package client

import (
	"context"
	"encoding/json"
	"fmt"
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
	RoundDelay      = 2 * time.Second
)

type NodeLatencyResponse struct {
	ServerNodeName string `json:"serverNodeName"`
	ServerPodName  string `json:"serverPodName"`
}

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

	res, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		log.Fatal(err)
	}
	end := time.Since(start)

	defer res.Body.Close()

	var resNodeInfo NodeLatencyResponse
	err = json.NewDecoder(res.Body).Decode(&resNodeInfo)
	if err != nil {
		panic(err)
	}

	log.Printf("trip completed: %s -> %s", c.nodeName, resNodeInfo.ServerNodeName)
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
		Observe(float64(end.Seconds()))
}

func (c *Client) startMetricsServer() {
	log.Println("starting client metrics server")

	prometheus.MustRegister(dnsLatencyHistogram)
	prometheus.MustRegister(connectionLatencyHistogram)
	prometheus.MustRegister(firstByteLatencyHistogram)
	prometheus.MustRegister(httpTotalLatencyHistogram)
	prometheus.MustRegister(httpTotalLatencySummary)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8081", nil))
}
