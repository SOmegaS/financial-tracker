package app

import (
	"context"
	"crypto/rsa"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"expensereader/internal/database"
	"expensereader/metrics"
	"expensereader/pkg/api"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	_ "github.com/lib/pq"
)

type App struct {
	api.UnimplementedApiServer
	db         *database.Queries
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

// NewApp создает экземпляр приложения expensereader
func NewApp(dbName, dbUser, dbHost, dbPort, dbPass string) (*App, error) {
	pubKeyData, err := os.ReadFile("secret/public.key")
	if err != nil {
		return nil, fmt.Errorf("failed to read public key: %w", err)
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(pubKeyData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?&sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName)
	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	db := database.New(dbConn)
	return &App{
		db:        db,
		publicKey: publicKey,
	}, nil
}

// GetReport возвращает агрегированный отчет по категориям
func (a *App) GetReport(ctx context.Context, req *api.GetReportRequest) (*api.GetReportResponse, error) {
	const method = "GetReport"
	start := time.Now()
	metrics.GRPCRequestsTotal.WithLabelValues(method).Inc()
	defer metrics.GRPCDurationSeconds.WithLabelValues(method).Observe(time.Since(start).Seconds())

	token, err := jwt.Parse(req.Jwt, func(token *jwt.Token) (interface{}, error) {
		return a.publicKey, nil
	})
	if err != nil {
		metrics.GRPCErrorsTotal.WithLabelValues(method).Inc()
		return nil, status.Errorf(codes.Unauthenticated, "parse token error: %v", err)
	}

	if token.Claims.Valid() != nil {
		metrics.GRPCErrorsTotal.WithLabelValues(method).Inc()
		return nil, status.Errorf(codes.Unauthenticated, "token is invalid: %v", err)
	}

	id, err := uuid.Parse(token.Claims.(jwt.MapClaims)["user_id"].(string))
	if err != nil {
		metrics.GRPCErrorsTotal.WithLabelValues(method).Inc()
		return nil, status.Errorf(codes.InvalidArgument, "invalid uuid")
	}
	log.Printf("Принят rpc запрос от пользоваля с user_id = %v", id)

	r, err := a.db.GetReport(ctx, id)
	if err != nil {
		metrics.GRPCErrorsTotal.WithLabelValues(method).Inc()
		return nil, status.Errorf(codes.Internal, "get report error: %v", err)
	}

	m := make(map[string]float64)
	for _, bill := range r {
		m[bill.Category] += bill.Amount
	}
	resp := &api.GetReportResponse{Report: m}
	return resp, nil
}

// GetBills возвращает список счетов по категории
func (a *App) GetBills(ctx context.Context, req *api.GetBillsRequest) (*api.GetBillsResponse, error) {
	const method = "GetBills"
	start := time.Now()
	metrics.GRPCRequestsTotal.WithLabelValues(method).Inc()
	defer metrics.GRPCDurationSeconds.WithLabelValues(method).Observe(time.Since(start).Seconds())

	token, err := jwt.Parse(req.Jwt, func(token *jwt.Token) (interface{}, error) {
		return a.publicKey, nil
	})
	if err != nil {
		metrics.GRPCErrorsTotal.WithLabelValues(method).Inc()
		return nil, status.Errorf(codes.Internal, "parse token error: %v", err)
	}
	if token.Claims.Valid() != nil {
		metrics.GRPCErrorsTotal.WithLabelValues(method).Inc()
		return nil, status.Errorf(codes.Unauthenticated, "token is invalid: %v", err)
	}

	id, err := uuid.Parse(token.Claims.(jwt.MapClaims)["user_id"].(string))
	if err != nil {
		metrics.GRPCErrorsTotal.WithLabelValues(method).Inc()
		return nil, status.Errorf(codes.InvalidArgument, "invalid uuid")
	}
	bills, err := a.db.GetBills(ctx, database.GetBillsParams{UserID: id, Category: req.Category})
	if err != nil {
		metrics.GRPCErrorsTotal.WithLabelValues(method).Inc()
		return nil, status.Errorf(codes.Internal, "get bills error: %v", err)
	}

	resp := &api.GetBillsResponse{Bills: make([]*api.Bill, 0, len(bills))}
	for _, bill := range bills {
		resp.Bills = append(resp.Bills, &api.Bill{
			Amount: bill.Amount,
			Name:   bill.Name,
			Ts:     timestamppb.New(bill.Tmstmp),
		})
	}
	return resp, nil
}
