package main

import (
	"net"
	"net/http"
	"os"

	"user-service/internal/app"
	"user-service/metrics"
	"user-service/pkg/api"

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
		sugar.Infow("Metrics endpoint listening on :2114/metrics")
		if err := http.ListenAndServe(":2114", nil); err != nil {
			sugar.Fatalw("metrics server failed", "error", err)
		}
	}()

	// Create app
	a, err := app.NewApp(sugar, dbName, dbUser, dbHost, dbPort, dbPass)
	if err != nil {
		sugar.Fatalw("failed to create app", "error", err)
	}

	// Инициализация приложения
	if err := a.Init(); err != nil {
		sugar.Fatalw("failed to initialize app", "error", err)
	}

	// Инициализация gRPC сервера
	listener, err := net.Listen("tcp", ":7777")
	if err != nil {
		sugar.Fatalw("failed to listen", "error", err)
	}
	defer listener.Close()

	s := grpc.NewServer()
	api.RegisterApiServer(s, a)

	// Запуск сервера
	sugar.Infow("User service started on :7777")
	if err := s.Serve(listener); err != nil {
		sugar.Fatalw("failed to serve", "error", err)
	}
}
