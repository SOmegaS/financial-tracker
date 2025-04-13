package app

import (
	"context"
	"crypto/rsa"
	"database/sql"
	"fmt"
	"os"
	"time"

	"expensereader/internal/database"
	"expensereader/pkg/api"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	_ "github.com/lib/pq"
)

type App struct {
	api.UnimplementedApiServer
	db         *database.Queries
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewApp(dbName, dbUser, dbPass string) (*App, error) {
	privKeyData, err := os.ReadFile("../../secret/private.key")
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %w", err)
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privKeyData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	// Optionally load public key (for verification elsewhere)
	pubKeyData, err := os.ReadFile("../../secret/public.key")
	if err != nil {
		return nil, fmt.Errorf("failed to read public key: %w", err)
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(pubKeyData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	connStr := fmt.Sprintf("postgres://%v:%v@localhost:5434/%v?&sslmode=disable", dbUser, dbPass, dbName)
	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	db := database.New(dbConn)
	return &App{
		db:         db,
		privateKey: privateKey,
		publicKey:  publicKey,
	}, nil
}

func (a *App) GetReport(ctx context.Context, req *api.GetReportRequest) (*api.GetReportResponse, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(req.Jwt, claims, func(token *jwt.Token) (interface{}, error) {
		return a.publicKey, nil
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "parse token error: %v", err)
	}

	if claims["exp"].(int64) > time.Now().Unix() {
		return nil, status.Errorf(codes.Unauthenticated, "token expired")
	}
	r, err := a.db.GetReport(ctx, claims["id"].(uuid.UUID))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get report error: %v", err)
	}
	m := make(map[string]float64)
	for _, bill := range r {
		m[bill.Category] += bill.Amount
	}
	resp := &api.GetReportResponse{
		Report: m,
	}
	return resp, nil
}
