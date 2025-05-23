package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	GRPCRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "expensereader_grpc_requests_total",
			Help: "Total number of gRPC requests received by expensereader",
		},
		[]string{"method"},
	)
	GRPCErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "expensereader_grpc_errors_total",
			Help: "Total number of errors in expensereader gRPC methods",
		},
		[]string{"method"},
	)
	GRPCDurationSeconds = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "expensereader_grpc_duration_seconds",
			Help:    "Histogram of gRPC method execution durations",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)
)

func Init() {
	prometheus.MustRegister(GRPCRequestsTotal, GRPCErrorsTotal, GRPCDurationSeconds)
}
