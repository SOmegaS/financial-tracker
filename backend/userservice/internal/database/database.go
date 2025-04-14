package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

const migrationsDir = "internal/database/migrations"

func Open(dbName, dbUser, dbPass string) (DBTX, error) {
	connStr := fmt.Sprintf("postgres://%v:%v@postgres-db-primary:5432/%v?&sslmode=disable", dbUser, dbPass, dbName)
	return sql.Open("postgres", connStr)
}

func RunMigrations(db DBTX, dbName string) (isChanged bool, err error) {
	driver, err := postgres.WithInstance(db.(*sql.DB), &postgres.Config{})
	if err != nil {
		return false, err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsDir,
		dbName,
		driver,
	)
	if err != nil {
		return false, err
	}
	err = m.Up()
	if err == nil {
		return true, nil
	}
	if errors.Is(err, migrate.ErrNoChange) {
		return false, nil
	}
	return false, err
}
