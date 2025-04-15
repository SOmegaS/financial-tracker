package main

import (
	"expensereader/internal/app"
	"expensereader/pkg/api"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

func main() {
	// Get env vars
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	a, err := app.NewApp(dbName, dbUser, dbHost, dbPort, dbPass)
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
	log.Println("Приложение запущено")
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
