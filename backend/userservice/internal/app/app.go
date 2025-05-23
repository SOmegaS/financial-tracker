package app

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"user-service/internal/database"
	"user-service/metrics"
	"user-service/pkg/api"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const expirationTime = time.Hour

type App struct {
	api.UnimplementedApiServer
	db         database.DBTX
	dbName     string
	queries    *database.Queries
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

// NewApp инициализирует подключение к базе данных
func NewApp(dbName, dbUser, dbHost, dbPort, dbPass string) (*App, error) {
	db, err := database.Open(dbName, dbUser, dbHost, dbPort, dbPass)
	if err != nil {
		return nil, err
	}
	return &App{
		db:     db,
		dbName: dbName,
	}, nil
}

// Init применяет миграции и загружает ключи
func (a *App) Init() error {
	isChanged, err := database.RunMigrations(a.db, a.dbName)
	if err != nil {
		return err
	}
	if isChanged {
		log.Println("Migrations applied")
	}
	a.queries = database.New(a.db)

	// Загрузка приватного ключа
	privKeyData, err := os.ReadFile("secret/private.key")
	if err != nil {
		return fmt.Errorf("failed to read private key: %w", err)
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privKeyData)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}
	a.privateKey = privateKey

	// Загрузка публичного ключа
	pubKeyData, err := os.ReadFile("secret/public.key")
	if err != nil {
		return fmt.Errorf("failed to read public key: %w", err)
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(pubKeyData)
	if err != nil {
		return fmt.Errorf("failed to parse public key: %w", err)
	}
	a.publicKey = publicKey

	return nil
}

// CreateSession создает JWT-токен сессии
func (a *App) CreateSession(id, userId string) (string, error) {
	claims := jwt.MapClaims{
		"id":      id,
		"exp":     time.Now().Add(expirationTime).Unix(),
		"user_id": userId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(a.privateKey)
}

// Register обрабатывает регистрацию нового пользователя
func (a *App) Register(ctx context.Context, req *api.RegisterRequest) (*api.RegisterResponse, error) {
	const method = "Register"
	start := time.Now()
	metrics.GRPCRequestsTotal.WithLabelValues(method).Inc()
	defer metrics.GRPCDurationSeconds.WithLabelValues(method).Observe(time.Since(start).Seconds())

	if len(req.Password) < 9 {
		metrics.GRPCErrorsTotal.WithLabelValues(method).Inc()
		return nil, status.Errorf(codes.InvalidArgument, "password must be longer than 8 characters")
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		metrics.GRPCErrorsTotal.WithLabelValues(method).Inc()
		return nil, status.Errorf(codes.Internal, "failed to hash password: %v", err)
	}

	userId, err := uuid.NewUUID()
	if err != nil {
		metrics.GRPCErrorsTotal.WithLabelValues(method).Inc()
		return nil, status.Errorf(codes.Internal, "failed to generate UUID")
	}

	err = a.queries.CreateUser(ctx, database.CreateUserParams{
		ID:       userId,
		Password: string(passHash),
		Username: req.GetUsername(),
	})
	if err != nil {
		metrics.GRPCErrorsTotal.WithLabelValues(method).Inc()
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	jwtToken, err := a.CreateSession(req.RequestId, userId.String())
	if err != nil {
		metrics.GRPCErrorsTotal.WithLabelValues(method).Inc()
		return nil, status.Errorf(codes.Internal, "failed to create session: %v", err)
	}

	return &api.RegisterResponse{
		Jwt: jwtToken,
	}, nil
}

// Login обрабатывает авторизацию пользователя
func (a *App) Login(ctx context.Context, req *api.LoginRequest) (*api.LoginResponse, error) {
	const method = "Login"
	start := time.Now()
	metrics.GRPCRequestsTotal.WithLabelValues(method).Inc()
	defer metrics.GRPCDurationSeconds.WithLabelValues(method).Observe(time.Since(start).Seconds())

	userInfo, err := a.queries.GetUserIdPassword(ctx, req.GetUsername())
	if err != nil {
		metrics.GRPCErrorsTotal.WithLabelValues(method).Inc()
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userInfo.Password), []byte(req.Password)); err != nil {
		metrics.GRPCErrorsTotal.WithLabelValues(method).Inc()
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, status.Errorf(codes.Unauthenticated, "invalid credentials")
		}
		return nil, status.Errorf(codes.Internal, "password comparison failed: %v", err)
	}

	jwtToken, err := a.CreateSession(req.RequestId, userInfo.ID.String())
	if err != nil {
		metrics.GRPCErrorsTotal.WithLabelValues(method).Inc()
		return nil, status.Errorf(codes.Internal, "failed to create session: %v", err)
	}

	return &api.LoginResponse{
		Jwt: jwtToken,
	}, nil
}
