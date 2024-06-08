package db

import (
	cfg "basic-trade/pkg/config"
	"context"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

const DSN = "postgres://%s:%s@%s:%d/%s?sslmode=%s"

func InitDB(d cfg.Database) (*pgxpool.Pool, error) {
	source := fmt.Sprintf(
		DSN,
		d.User,
		d.Password,
		d.Host,
		d.Port,
		d.Name,
		d.SslMode,
	)
	connPool, err := pgxpool.New(context.Background(), source)

	if err != nil {
		return nil, err
	}

	migration, err := migrate.New(d.MigrationURL, source)
	if err != nil {
		log.Fatal("Couldn't cretae migration: ", err)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Couldn't run migration up: ", err)
	}

	return connPool, nil
}
