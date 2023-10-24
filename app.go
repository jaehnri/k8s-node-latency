package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const EnvKubePodName = "KUBE_POD_NAME"

type Server struct {
	podName string
}

func NewServer() *Server {
	return &Server{
		podName: os.Getenv(EnvKubePodName),
	}
}

func (s *Server) handlePing(w http.ResponseWriter, r *http.Request) {
	pingCounter.Inc()
	fmt.Fprint(w, "pong")
}

func (s *Server) Run() {
	// Register Prometheus metrics
	prometheus.MustRegister(pingCounter)

	// Start HTTP server
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/ping", s.handlePing)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
