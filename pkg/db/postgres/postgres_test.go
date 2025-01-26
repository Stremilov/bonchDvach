package postgres_test

import (
	"bonchDvach/pkg/db/postgres"
	"testing"

	"github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
)

func TestMigrate(t *testing.T) {
	pool, err := postgres.New("host=localhost user=postgres password=postgres dbname=bonchdvach sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	db := stdlib.OpenDB(*pool.Config().ConnConfig)
	defer db.Close()
	err = goose.Down(db, "migrations")
	if err != nil {
		t.Fatal(err)
	}
}
