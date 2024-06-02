package db

import (
	cfg "basic-trade/pkg/config"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

const DSN = "host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s"

func InitDB(d cfg.Database) (*pgxpool.Pool, error) {
	connPool, err := pgxpool.New(
		context.Background(),
		fmt.Sprintf(
			DSN,
			d.Host,
			d.User,
			d.Password,
			d.Name,
			d.Port,
			d.SslMode,
			d.Timezone,
		),
	)

	if err != nil {
		return nil, err
	}

	return connPool, nil
}
