package client

import "github.com/prometheus/client_golang/prometheus"

var (
	dnsLatencyHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "http_dns_request_duration_ms",
		Help:    "Histogram of DNS requests latency",
		Buckets: []float64{1, 2, 4, 8, 16, 32, 64, 128, 256, 512},
	})

	connectionLatencyHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "http_connection_request_duration_ms",
		Help:    "Histogram of connection establishment latency",
		Buckets: []float64{1, 2, 4, 8, 16, 32, 64, 128, 256, 512},
	})

	firstByteLatencyHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "http_first_byte_request_duration_ms",
		Help:    "Histogram of latency of response first byte",
		Buckets: []float64{1, 2, 4, 8, 16, 32, 64, 128, 256, 512},
	})

	httpTotalLatencyHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "http_total_ping_request_duration_ms",
		Help:    "Histogram of HTTP ping requests latency",
		Buckets: []float64{1, 2, 4, 8, 16, 32, 64, 128, 256, 512},
	})
)
