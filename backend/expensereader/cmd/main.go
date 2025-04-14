package main

import (
	"expensereader/internal/app"
	"expensereader/pkg/api"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	// Get env vars
	dbUser := "ivang"
	dbPass := "ivang"
	dbName := "db"
	a, err := app.NewApp(dbName, dbUser, dbPass)
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
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
