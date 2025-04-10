package main

import (
	"log"
	"os"

	"financial-tracker/internal/app"
)

func main() {
	dbUser := os.Getenv("POSTGRES_USER")
	dbPass := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	a, err := app.NewApp(
		dbName,
		dbUser,
		dbPass,
	)
	if err != nil {
		log.Fatal(err)
	}
	a.Init()
	a.Run()
}
