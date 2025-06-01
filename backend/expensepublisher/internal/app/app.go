package app

import (
	"context"
	"crypto/rsa"
	"errors"
	"expensepublisher/pkg/api"
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"expensepublisher/metrics"
)

type App struct {
	api.UnimplementedApiServer
	writer    *kafka.Writer
	publicKey *rsa.PublicKey
	logger    *zap.SugaredLogger
}

func NewApp(logger *zap.SugaredLogger, kafkaHostPort, topicName string) (*App, error) {
	pubKeyStr := os.Getenv("PUBLIC_KEY")
	if pubKeyStr == "" {
		logger.Error("PUBLIC_KEY environment variable is not set")
		return nil, fmt.Errorf("PUBLIC_KEY environment variable is not set")
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pubKeyStr))
	if err != nil {
		logger.Errorw("failed to parse public key", "error", err)
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	logger.Info(publicKey)
	logger.Info(topicName)
	logger.Info(kafkaHostPort)
	w := &kafka.Writer{
		Addr:         kafka.TCP(kafkaHostPort),
		Topic:        topicName,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
	}

	return &App{
		writer:    w,
		publicKey: publicKey,
		logger:    logger,
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
	a.logger.Info("GOT MESSAGE CREATE BILL")
	start := time.Now()
	metrics.RequestsTotal.WithLabelValues("CreateBill").Inc()
	defer metrics.RequestDuration.WithLabelValues("CreateBill").Observe(time.Since(start).Seconds())

	token, err := jwt.Parse(msg.Jwt, func(token *jwt.Token) (interface{}, error) {
		return a.publicKey, nil
	})
	if err != nil {
		metrics.ErrorsTotal.WithLabelValues("CreateBill").Inc()
		a.logger.Errorw("token parse error", "error", err)
		return nil, status.Errorf(codes.Unauthenticated, "parse token error: %v", err)
	}

	if token.Claims.Valid() != nil {
		metrics.ErrorsTotal.WithLabelValues("CreateBill").Inc()
		a.logger.Errorw("token invalid", "error", err)
		return nil, status.Errorf(codes.Unauthenticated, "token is invalid: %v", err)
	}

	id, err := uuid.Parse(token.Claims.(jwt.MapClaims)["user_id"].(string))
	if err != nil {
		metrics.ErrorsTotal.WithLabelValues("CreateBill").Inc()
		a.logger.Errorw("invalid uuid", "error", err)
		return nil, status.Errorf(codes.InvalidArgument, "invalid uuid")
	}

	a.logger.Infow("Принят rpc запрос от пользователя", "userID", id)

	if err := validateCreateBillMessage(msg); err != nil {
		metrics.ErrorsTotal.WithLabelValues("CreateBill").Inc()
		a.logger.Errorw("validation error", "error", err)
		return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}

	err = a.publishMessage(ctx, id, msg)
	if err != nil {
		metrics.ErrorsTotal.WithLabelValues("CreateBill").Inc()
		a.logger.Errorw("failed to publish message", "error", err)
		return nil, status.Errorf(codes.Internal, "internal error: %v", err)
	}

	a.logger.Infow("Successfully published CreateBill message", "userID", id)
	return &emptypb.Empty{}, nil
}

func (a *App) publishMessage(ctx context.Context, userId uuid.UUID, msg *api.BillMessage) error {
	a.logger.Info(msg)
	writeMessage := &api.CreateBillMessage{
		Name:      msg.Name,
		Amount:    msg.Amount,
		Category:  msg.Category,
		Timestamp: msg.Timestamp,
		UserId:    userId.String(),
	}
	bytes, err := proto.Marshal(writeMessage)
	if err != nil {
		a.logger.Errorw("failed to marshal message", "error", err)
		return err
	}
	err = a.writer.WriteMessages(ctx, kafka.Message{
		Value: bytes,
	})
	if err != nil {
		a.logger.Errorw("failed to write to kafka", "error", err)
	}
	return err
}
