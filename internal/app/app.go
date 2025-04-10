package app

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"financial-tracker/internal/database"
	"financial-tracker/pkg/api"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const ExpirationTime = time.Hour

type App struct {
	api.UnimplementedApiServer
	db        database.DBTX
	dbName    string
	queries   *database.Queries
	secretKey []byte
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
	_, err = rand.Read(a.secretKey)
	return err
}

func (a *App) Run() {
	fmt.Print("Hello, world!")
}

func (a *App) CreateSession(id string) (string, error) {
	claims := jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(ExpirationTime).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.secretKey)
}

func (a *App) Register(ctx context.Context, req *api.RegisterRequest) (*api.RegisterResponse, error) {
	passHash, err := bcrypt.GenerateFromPassword([]byte(req.Pass), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate password hash: %v", err)
	}
	err = a.queries.CreateUser(ctx, database.CreateUserParams{
		ID:       req.Id,
		PassHash: sql.NullString{String: string(passHash), Valid: true},
		Username: sql.NullString{String: req.Username, Valid: true},
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}
	jwt, err := a.CreateSession(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session: %v", err)
	}
	return &api.RegisterResponse{
		Jwt: jwt,
	}, nil
}

func (a *App) DeleteUser(context.Context, *api.DeleteUserRequest) (*api.DeleteUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUser not implemented")
}

func (a *App) Login(ctx context.Context, req *api.LoginRequest) (*api.LoginResponse, error) {
	user, err := a.queries.GetUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if !user.PassHash.Valid {
		return nil, status.Errorf(codes.DataLoss, "password corrupted: %v", err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PassHash.String), []byte(req.Pass)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, status.Errorf(codes.Unauthenticated, "invalid password: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to compare password hash: %v", err)
	}
	jwt, err := a.CreateSession(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session: %v", err)
	}
	return &api.LoginResponse{
		Jwt: jwt,
	}, nil
}

func (a *App) Logout(context.Context, *api.LogoutRequest) (*api.LogoutResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Logout not implemented")
}

func (a *App) CreateCategory(context.Context, *api.CreateCategoryRequest) (*api.CreateCategoryResponse, error) {

	return nil, status.Errorf(codes.Unimplemented, "method CreateCategory not implemented")
}

func (a *App) DeleteCategory(context.Context, *api.DeleteCategoryRequest) (*api.DeleteCategoryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteCategory not implemented")
}

func (a *App) ListCategories(context.Context, *api.ListCategoriesRequest) (*api.ListCategoriesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListCategories not implemented")
}

func (a *App) CreateBill(context.Context, *api.CreateBillRequest) (*api.CreateBillResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateBill not implemented")
}

func (a *App) DeleteBill(context.Context, *api.DeleteBillRequest) (*api.DeleteBillResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteBill not implemented")
}

func (a *App) ListBills(context.Context, *api.ListBillsRequest) (*api.ListBillsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListBills not implemented")
}
