package main

import (
	"log"
	"net"
	"os"
	"user-service/internal/app"
	"user-service/pkg/api"

	"google.golang.org/grpc"
)

func main() {
	// Get env vars
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Create app
	a, err := app.NewApp(
		dbName,
		dbUser,
		dbHost,
		dbPort,
		dbPass,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Инициализация gRPC сервера
	listener, err := net.Listen("tcp", ":7777")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer listener.Close()

	s := grpc.NewServer()
	api.RegisterApiServer(s, a)

	// Инициализация приложения
	if err := a.Init(); err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	a.Run()

	// Запуск сервера
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
