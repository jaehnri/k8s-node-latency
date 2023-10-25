package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	tcpPingCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "tcp_ping_requests_total",
		Help: "Total number of TCP ping requests",
	})

	httpPingCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "http_ping_requests_total",
		Help: "Total number of HTTP ping requests",
	})
)
