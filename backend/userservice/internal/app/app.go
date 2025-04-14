package app

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"os"
	"time"

	"user-service/internal/database"
	"user-service/pkg/api"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const ExpirationTime = time.Hour

type App struct {
	api.UnimplementedApiServer
	db         database.DBTX
	dbName     string
	queries    *database.Queries
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func NewApp(dbName, dbUser, dbPass string) (*App, error) {
	db, err := database.Open(dbName, dbUser, dbPass)
	if err != nil {
		return nil, err
	}
	return &App{
		db:     db,
		dbName: dbName,
	}, nil
}

func (a *App) Init() error {
	isChanged, err := database.RunMigrations(a.db, a.dbName)
	if err != nil {
		return err
	}
	if isChanged {
		log.Println("Migrations applied")
	}
	a.queries = database.New(a.db)

	privKeyData, err := os.ReadFile("secret/private.key")
	if err != nil {
		return fmt.Errorf("failed to read private key: %w", err)
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privKeyData)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}
	a.privateKey = privateKey

	// Optionally load public key (for verification elsewhere)
	pubKeyData, err := os.ReadFile("secret/public.key")
	if err != nil {
		return fmt.Errorf("failed to read public key: %w", err)
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(pubKeyData)
	if err != nil {
		return fmt.Errorf("failed to parse public key: %w", err)
	}
	a.publicKey = publicKey

	return err
}

func (a *App) Run() {
	fmt.Print("Hello, world!")
}

func (a *App) CreateSession(id string, userId string) (string, error) {
	claims := jwt.MapClaims{
		"id":      id,
		"exp":     time.Now().Add(ExpirationTime).Unix(),
		"user_id": userId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(a.privateKey)
}

func (a *App) Register(ctx context.Context, req *api.RegisterRequest) (*api.RegisterResponse, error) {
	passHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate password hash: %v", err)
	}
	userId, err := uuid.NewUUID()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate uuid")
	}
	err = a.queries.CreateUser(ctx, database.CreateUserParams{
		ID:       userId,
		Password: string(passHash),
		Username: req.GetUsername(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}
	jwt, err := a.CreateSession(req.RequestId, userId.String())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session: %v", err)
	}
	return &api.RegisterResponse{
		Jwt: jwt,
	}, nil
}

func (a *App) Login(ctx context.Context, req *api.LoginRequest) (*api.LoginResponse, error) {
	userInfo, err := a.queries.GetUserIdPassword(ctx, req.GetUsername())
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(userInfo.Password), []byte(req.Pass)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, status.Errorf(codes.Unauthenticated, "invalid password: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to compare password hash: %v", err)
	}
	jwt, err := a.CreateSession(req.RequestId, userInfo.ID.String())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session: %v", err)
	}
	return &api.LoginResponse{
		Jwt: jwt,
	}, nil
}
