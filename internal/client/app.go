package client

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"net/http"
	"net/http/httptrace"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Client struct {
	podName        string
	metricsAddress string

	kubeClient *kubernetes.Clientset
	dest       *v1.Service
}

func NewClient() *Client {
	serviceName := "your-service-name"
	namespace := "your-namespace"

	config, err := rest.InClusterConfig()
	if err != nil {
		log.Panic("couldn't fetch the in-cluster config", err)
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panic("couldn't create kubernetes config")
	}

	service, err := kubeClient.CoreV1().Services(namespace).Get(context.Background(), serviceName, metav1.GetOptions{})
	if err != nil {
		log.Panic("couldn't get server ClusterIP")
	}

	clusterIP := service.Spec.ClusterIP
	fmt.Printf("ClusterIP of service %s in namespace %s is %s\n", serviceName, namespace, clusterIP)

	return &Client{
		kubeClient: kubeClient,
	}
}

func (c *Client) testLoop() {
	for {
		c.sendHTTPPing()
		log.Println("round over!")

		time.Sleep(2 * time.Second)
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
	//url := fmt.Sprintf("http://%s:8080/ping", c.dest.Spec.ClusterIP)
	url := fmt.Sprintf("http://%s:8080/ping", "localhost")
	req, _ := http.NewRequest("GET", url, nil)

	var start, connect, dns time.Time

	trace := &httptrace.ClientTrace{
		DNSStart: func(dsi httptrace.DNSStartInfo) { dns = time.Now() },
		DNSDone: func(ddi httptrace.DNSDoneInfo) {
			dnsLatencyHistogram.Observe(float64(time.Since(dns).Milliseconds()))
		},

		ConnectStart: func(network, addr string) { connect = time.Now() },
		ConnectDone: func(network, addr string, err error) {
			connectionLatencyHistogram.Observe(float64(time.Since(connect).Milliseconds()))
		},

		GotFirstResponseByte: func() {
			firstByteLatencyHistogram.Observe(float64(time.Since(start).Milliseconds()))
		},
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	start = time.Now()
	if _, err := http.DefaultTransport.RoundTrip(req); err != nil {
		log.Fatal(err)
	}
	end := time.Since(start)

	httpTotalLatencyHistogram.Observe(float64(end.Milliseconds()))
}

func (c *Client) startMetricsServer() {
	log.Println("starting client metrics server")

	prometheus.MustRegister(dnsLatencyHistogram)
	prometheus.MustRegister(connectionLatencyHistogram)
	prometheus.MustRegister(firstByteLatencyHistogram)
	prometheus.MustRegister(httpTotalLatencyHistogram)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8081", nil))
}
