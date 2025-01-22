package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Postgres struct {
	Pool *pgxpool.Pool
}

func New(connString string) (*Postgres, error) {
	ctx := context.Background()

	config, _ := pgxpool.ParseConfig(connString)
	config.MaxConns = 50

	pool, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to db: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot ping db: %w", err)
	}

	return &Postgres{Pool: pool}, nil
}

func applyMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	return nil
}

//Методы для работы с БД
