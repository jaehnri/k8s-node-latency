package server

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const EnvKubePodName = "KUBE_POD_NAME"
const EnvKubeNodeName = "KUBE_NODE_NAME"

type NodeLatencyResponse struct {
	ServerNodeName string `json:"serverNodeName"`
	ServerPodName  string `json:"serverPodName"`
}

type Server struct {
	podName  string
	nodeName string

	tcpAddress  string
	httpAddress string
}

func NewServer(tcpAddress, httpAddress string) *Server {
	podName, exists := os.LookupEnv(EnvKubePodName)
	if !exists {
		log.Panic("couldn't retrieve podname")
	}

	nodeName, exists := os.LookupEnv(EnvKubeNodeName)
	if !exists {
		log.Panic("couldn't retrieve node name")
	}
	return &Server{
		nodeName:    nodeName,
		podName:     podName,
		tcpAddress:  tcpAddress,
		httpAddress: httpAddress,
	}
}

func (s *Server) generateResponse() []byte {
	serverInfo := NodeLatencyResponse{
		ServerNodeName: s.nodeName,
		ServerPodName:  s.podName,
	}

	responseJSON, err := json.Marshal(serverInfo)
	if err != nil {
		log.Println("couldn't marshal response to JSON", err)
	}

	return responseJSON
}

func (s *Server) handleHTTPPing(w http.ResponseWriter, r *http.Request) {
	log.Println("received call to HTTP /ping")
	httpPingCounter.Inc()
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	responseJSON := s.generateResponse()

	w.Write(responseJSON)
}

func (s *Server) handleTCPPing(conn net.Conn) {
	log.Println("received call to TCP /ping")
	tcpPingCounter.Inc()

	responseJSON := s.generateResponse()

	_, err := conn.Write(responseJSON)
	if err != nil {
		log.Println("couldn't respond TCP request", err)
	}
	_ = conn.Close()
}

func (s *Server) Run() {
	prometheus.MustRegister(tcpPingCounter)
	prometheus.MustRegister(httpPingCounter)
	http.Handle("/metrics", promhttp.Handler())

	go s.startTCPServer(s.handleTCPPing)
	go s.startHTTPServer(s.handleHTTPPing)

	sigint := make(chan os.Signal, 1)
	// interrupt signal sent from terminal
	signal.Notify(sigint, os.Interrupt)
	// sigterm signal sent from kubernetes
	signal.Notify(sigint, syscall.SIGTERM)

	<-sigint
	log.Fatal("SIGINT received, shutting down server")
}

func (s *Server) startTCPServer(handler func(conn net.Conn)) {
	log.Println("starting ping TCP server")

	listen, err := net.Listen("tcp", s.tcpAddress)
	if err != nil {
		log.Panic("failed to bind TCP server", err)
	}
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Panic("couldn't accept connection: ", err)
		}
		go handler(conn)
	}
}

func (s *Server) startHTTPServer(handler func(w http.ResponseWriter, r *http.Request)) {
	log.Println("starting ping HTTP server")

	http.HandleFunc("/ping", handler)
	log.Fatal(http.ListenAndServe(s.httpAddress, nil))
}
