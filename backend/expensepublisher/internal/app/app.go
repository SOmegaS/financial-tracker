package app

import (
	"context"
	"errors"
	"expensepublisher/pkg/api"
	"fmt"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"math"
	"strings"
)

type App struct {
	api.UnimplementedApiServer
	writer *kafka.Writer
}

func NewApp() (*App, error) {
	w := &kafka.Writer{
		Addr:         kafka.TCP("kafka-moscow:9092"),
		Topic:        "write-bills",
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
	}
	return &App{writer: w}, nil
}

func validateCreateBillMessage(msg *api.CreateBillMessage) error {
	if strings.TrimSpace(msg.Name) == "" {
		return errors.New("name is required")
	}
	if strings.TrimSpace(msg.Category) == "" {
		return errors.New("category is required")
	}
	if strings.TrimSpace(msg.UserId) == "" {
		return errors.New("user_id is required")
	}
	if _, err := uuid.Parse(msg.UserId); err != nil {
		return fmt.Errorf("user_id must be a valid UUID: %v", err)
	}
	if msg.Amount <= 0 || msg.Amount >= 1_000_000_000_000 {
		return errors.New("amount must be greater than 0 and less than 1_000_000_000_000")
	}
	if math.IsNaN(msg.Amount) || math.IsInf(msg.Amount, 0) {
		return errors.New("amount must be a finite number")
	}
	if msg.Timestamp == nil || msg.Timestamp.AsTime().IsZero() {
		return errors.New("timestamp is required and must be valid")
	}
	return nil
}

func (a *App) CreateBill(ctx context.Context, msg *api.CreateBillMessage) (*emptypb.Empty, error) {
	log.Println("Принят rpc вызов")

	if err := validateCreateBillMessage(msg); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}
	err := a.CreateBillPublisher(ctx, msg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (a *App) CreateBillPublisher(ctx context.Context, msg *api.CreateBillMessage) error {
	bytes, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	err = a.writer.WriteMessages(ctx, kafka.Message{
		Value: bytes,
	})
	return err
}
