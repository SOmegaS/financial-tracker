package main

import (
	"expensepublisher/internal/app"
	"expensepublisher/pkg/api"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	// Create app
	a, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.Listen("tcp", ":7777")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer listener.Close()
	s := grpc.NewServer()
	api.RegisterApiServer(s, a)

	log.Println("Успешная настройка")

	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
