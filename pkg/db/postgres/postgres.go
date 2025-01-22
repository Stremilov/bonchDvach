package postgres

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
)

type Postgres struct {
	Pool *pgxpool.Pool
}

const (
	pathToMigrations = "./migrations"
)

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

	err = applyMigrations(pool)
	if err != nil {
		return nil, fmt.Errorf("cannot apply migrations: %w", err)
	}

	return &Postgres{Pool: pool}, nil
}

func applyMigrations(pool *pgxpool.Pool) error {
	db := stdlib.OpenDB(*pool.Config().ConnConfig)
	defer db.Close()

	migrationsDir := pathToMigrations

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %v", err)
	}

	if err := goose.Up(db, migrationsDir); err != nil {
		return fmt.Errorf("failed to apply migrations: %v", err)
	}

	log.Println("Migrations applied successfully!")

	return nil
}

//Методы для работы с БД
