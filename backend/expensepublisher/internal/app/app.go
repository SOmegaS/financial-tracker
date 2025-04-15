package app

import (
	"context"
	"crypto/rsa"
	"errors"
	"expensepublisher/pkg/api"
	"fmt"
	"log"
	"math"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type App struct {
	api.UnimplementedApiServer
	writer    *kafka.Writer
	publicKey *rsa.PublicKey
}

func NewApp(kafkaHostPort, topicName string) (*App, error) {
	pubKeyData, err := os.ReadFile("secret/public.key")
	if err != nil {
		return nil, fmt.Errorf("failed to read public key: %w", err)
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(pubKeyData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	w := &kafka.Writer{
		Addr:         kafka.TCP(kafkaHostPort),
		Topic:        topicName,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
	}

	return &App{
		writer:    w,
		publicKey: publicKey,
	}, nil
}

func validateCreateBillMessage(msg *api.BillMessage) error {
	if strings.TrimSpace(msg.Name) == "" {
		return errors.New("name is required")
	}
	if strings.TrimSpace(msg.Category) == "" {
		return errors.New("category is required")
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

func (a *App) CreateBill(ctx context.Context, msg *api.BillMessage) (*emptypb.Empty, error) {

	token, err := jwt.Parse(msg.Jwt, func(token *jwt.Token) (interface{}, error) {
		return a.publicKey, nil
	})
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "parse token error: %v", err)
	}

	if token.Claims.Valid() != nil {
		return nil, status.Errorf(codes.Unauthenticated, "token is invalid: %v", err)
	}

	id, err := uuid.Parse(token.Claims.(jwt.MapClaims)["user_id"].(string))
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid uuid")
	}
	log.Printf("Принят rpc запрос от пользователя с id = %v", id)

	if err := validateCreateBillMessage(msg); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}

	err = a.publishMessage(ctx, id, msg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (a *App) publishMessage(ctx context.Context, userId uuid.UUID, msg *api.BillMessage) error {
	writeMessage := &api.CreateBillMessage{
		Name:      msg.Name,
		Amount:    msg.Amount,
		Category:  msg.Category,
		Timestamp: msg.Timestamp,
		UserId:    userId.String(),
	}
	bytes, err := proto.Marshal(writeMessage)
	if err != nil {
		return err
	}
	err = a.writer.WriteMessages(ctx, kafka.Message{
		Value: bytes,
	})
	return err
}
