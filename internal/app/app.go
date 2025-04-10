package app

import (
	"financial-tracker/internal/database"
	"fmt"
	"log"
)

type App struct {
	db     database.DBTX
	dbName string
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
	return nil
}

func (a *App) Run() {
	fmt.Print("Hello, world!")
}
