package repositories

import (
	"bonchDvach/pkg/handlers"
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type UserStorage struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) handlers.UserRepository {
	return UserStorage{
		db: db,
	}
}

// CreateUser implements db.UserRepository.
func (r UserStorage) CreateUser(ctx context.Context, userIP string) error {
	query := "INSERT INTO users (ip) VALUES ($1)"

	_, err := r.db.Exec(ctx, query, userIP)

	return err
}

// GetUser implements db.UserRepository.
func (r UserStorage) GetUser() {
	panic("unimplemented")
}
