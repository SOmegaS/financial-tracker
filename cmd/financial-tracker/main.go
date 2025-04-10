package main

import (
	"log"
	"net"
	"os"

	"financial-tracker/internal/app"
	"financial-tracker/pkg/api"

	"google.golang.org/grpc"
)

func main() {
	// Get env vars
	dbUser := os.Getenv("POSTGRES_USER")
	dbPass := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")

	// Create app
	a, err := app.NewApp(
		dbName,
		dbUser,
		dbPass,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Register gRPC server
	listener, err := net.Listen("tcp", ":7777")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer listener.Close()
	s := grpc.NewServer()
	api.RegisterApiServer(s, a)

	// Init app
	a.Init()

	// Serve
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
