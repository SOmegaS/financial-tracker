package main

import (
	"expensepublisher/internal/app"
	"expensepublisher/pkg/api"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func main() {
	kafkaHostPort := os.Getenv("KAFKA_HOST_PORT")
	topicName := os.Getenv("TOPIC_NAME")
	a, err := app.NewApp(kafkaHostPort, topicName)
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

	log.Println("Приложение запущено")

	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
