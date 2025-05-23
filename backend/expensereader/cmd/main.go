package main

import (
	"expensereader/internal/app"
	"expensereader/metrics"
	"log"
	"net"
	"net/http"
	"os"

	"expensereader/pkg/api"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

func main() {
	// Инициализация метрик
	metrics.Init()

	// HTTP для Prometheus
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("Metrics endpoint listening on :2113/metrics")
		log.Fatal(http.ListenAndServe(":2113", nil))
	}()

	// gRPC
	// Get env vars...
	a, err := app.NewApp(
		os.Getenv("DB_NAME"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_PASS"),
	)
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}

	listener, err := net.Listen("tcp", ":7777")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer listener.Close()
	s := grpc.NewServer()
	api.RegisterApiServer(s, a)

	log.Println("gRPC server listening on :7777")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
