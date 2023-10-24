package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	pingCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ping_requests_total",
		Help: "Total number of ping requests",
	})
)
