package client

import "github.com/prometheus/client_golang/prometheus"

var (
	dnsLatencyHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_dns_request_duration_ms",
			Help:    "Histogram of DNS requests latency",
			Buckets: []float64{1, 2, 4, 8, 16, 32, 64, 128, 256, 512},
		},
		[]string{"clientPodName", "clientNodeName", "serverPodName", "serverNodeName"},
	)

	connectionLatencyHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_connection_request_duration_ms",
			Help:    "Histogram of connection establishment latency",
			Buckets: []float64{1, 2, 4, 8, 16, 32, 64, 128, 256, 512},
		},
		[]string{"clientPodName", "clientNodeName", "serverPodName", "serverNodeName"},
	)

	firstByteLatencyHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_first_byte_request_duration_ms",
			Help:    "Histogram of latency of response first byte",
			Buckets: []float64{1, 2, 4, 8, 16, 32, 64, 128, 256, 512},
		},
		[]string{"clientPodName", "clientNodeName", "serverPodName", "serverNodeName"},
	)

	httpTotalLatencyHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_total_ping_request_duration_ms",
			Help:    "Histogram of HTTP ping requests latency",
			Buckets: []float64{1, 2, 4, 8, 16, 32, 64, 128, 256, 512},
		},
		[]string{"clientPodName", "clientNodeName", "serverPodName", "serverNodeName"},
	)

	httpTotalLatencySummary = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "http_total_latency_durations",
			Help:       "Ping durations in ms",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"clientPodName", "clientNodeName", "serverPodName", "serverNodeName"},
	)

	httpOneWayTripLatencySummary = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "http_one_way_trip_latency_durations",
			Help:       "One durations in ms",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"clientPodName", "clientNodeName", "serverPodName", "serverNodeName"},
	)
)
