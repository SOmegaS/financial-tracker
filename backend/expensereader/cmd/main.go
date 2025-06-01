package main

import (
	"expensereader/internal/app"
	"expensereader/metrics"
	"net"
	"net/http"
	"os"

	"expensereader/pkg/api"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	// Инициализация метрик
	metrics.Init()
	metrics.GRPCRequestsTotal.WithLabelValues("GetReport").Inc()

	// HTTP для Prometheus
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		sugar.Infow("Metrics endpoint listening on", "address", ":2113/metrics")
		if err := http.ListenAndServe(":2113", nil); err != nil {
			sugar.Fatalw("metrics server failed", "error", err)
		}
	}()

	// gRPC
	// Get env vars...
	a, err := app.NewApp(
		sugar,
		os.Getenv("DB_NAME"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_PASS"),
	)
	if err != nil {
		sugar.Fatalw("failed to create app", "error", err)
	}

	listener, err := net.Listen("tcp", ":7777")
	if err != nil {
		sugar.Fatalw("failed to listen", "error", err)
	}
	defer listener.Close()
	s := grpc.NewServer()
	api.RegisterApiServer(s, a)

	sugar.Infow("gRPC server listening on", "address", ":7777")
	if err := s.Serve(listener); err != nil {
		sugar.Fatalw("failed to serve", "error", err)
	}
}
