package database

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const migrationsDir = "internal/database/migrations"

func RunMigrations(dbName, dbUser, dbPass string) (isChanged bool, err error) {
	m, err := migrate.New(
		"file://"+migrationsDir,
		fmt.Sprintf("postgres://%v:%v@localhost:5432/%v?&sslmode=disable", dbUser, dbPass, dbName),
	)
	if err != nil {
		return false, err
	}
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
