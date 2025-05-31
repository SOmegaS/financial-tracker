// main.go
package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"user-service/internal/app"
	"user-service/metrics"
	"user-service/pkg/api"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

func main() {
	// Get env vars
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Инициализация метрик и HTTP endpoint
	metrics.Init()
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("Metrics endpoint listening on :2114/metrics")
		log.Fatal(http.ListenAndServe(":2114", nil))
	}()

	// Create app
	a, err := app.NewApp(dbName, dbUser, dbHost, dbPort, dbPass)
	if err != nil {
		log.Fatal(err)
	}

	// Инициализация приложения
	if err := a.Init(); err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	// Инициализация gRPC сервера
	listener, err := net.Listen("tcp", ":7777")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer listener.Close()

	s := grpc.NewServer()
	api.RegisterApiServer(s, a)

	// Запуск сервера
	log.Println("User service started on :7777")
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
