package main

import (
	"expensepublisher/internal/app"
	"expensepublisher/metrics"
	"expensepublisher/pkg/api"
	"net"
	"net/http"
	"os"

	"go.uber.org/zap"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	kafkaHostPort := os.Getenv("KAFKA_HOST_PORT")
	topicName := os.Getenv("TOPIC_NAME")

	metrics.InitMetrics()
	metrics.RequestsTotal.WithLabelValues("CreateBill").Inc()

	a, err := app.NewApp(sugar, kafkaHostPort, topicName)
	if err != nil {
		sugar.Fatalw("failed to create app", "error", err)
	}

	// HTTP server for Prometheus metrics
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		sugar.Infow("Метрики Prometheus доступны по адресу", "address", ":2112/metrics")
		if err := http.ListenAndServe(":2112", nil); err != nil {
			sugar.Fatalw("metrics server failed", "error", err)
		}
	}()

	listener, err := net.Listen("tcp", ":7777")
	if err != nil {
		sugar.Fatalw("failed to listen", "error", err)
	}
	defer listener.Close()

	s := grpc.NewServer()
	api.RegisterApiServer(s, a)

	sugar.Infow("Приложение запущено", "port", 7777)

	if err := s.Serve(listener); err != nil {
		sugar.Fatalw("failed to serve", "error", err)
	}
}
