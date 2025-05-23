package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	RequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "expense_publisher_requests_total",
			Help: "Total number of gRPC requests received",
		},
		[]string{"method"},
	)

	ErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "expense_publisher_errors_total",
			Help: "Total number of errors occurred",
		},
		[]string{"method"},
	)

	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "expense_publisher_request_duration_seconds",
			Help:    "Histogram of response time for handler",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)
)

func InitMetrics() {
	prometheus.MustRegister(RequestsTotal, ErrorsTotal, RequestDuration)
	RequestsTotal.WithLabelValues("CreateBill")
	ErrorsTotal.WithLabelValues("CreateBill")
	RequestDuration.WithLabelValues("CreateBill")
}
