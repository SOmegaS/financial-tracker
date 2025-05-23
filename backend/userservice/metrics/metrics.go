package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	GRPCRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "userservice_grpc_requests_total",
			Help: "Total number of gRPC requests received",
		},
		[]string{"method"},
	)
	GRPCErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "userservice_grpc_errors_total",
			Help: "Total number of errors in gRPC methods",
		},
		[]string{"method"},
	)
	GRPCDurationSeconds = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "userservice_grpc_duration_seconds",
			Help:    "Histogram of response time for gRPC methods",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)
)

func Init() {
	prometheus.MustRegister(GRPCRequestsTotal, GRPCErrorsTotal, GRPCDurationSeconds)
}
