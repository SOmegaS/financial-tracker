package main

import (
	"expensepublisher/internal/app"
	"expensepublisher/metrics"
	"expensepublisher/pkg/api"
	"log"
	"net"
	"net/http"
	"os"

	"google.golang.org/grpc"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	kafkaHostPort := os.Getenv("KAFKA_HOST_PORT")
	topicName := os.Getenv("TOPIC_NAME")

	metrics.InitMetrics()

	a, err := app.NewApp(kafkaHostPort, topicName)
	if err != nil {
		log.Fatal(err)
	}

	// HTTP server for Prometheus metrics
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("Метрики Prometheus доступны по адресу :2112/metrics")
		log.Fatal(http.ListenAndServe(":2112", nil))
	}()

	listener, err := net.Listen("tcp", ":7777")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer listener.Close()

	s := grpc.NewServer()
	api.RegisterApiServer(s, a)

	log.Println("Приложение запущено")

	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
