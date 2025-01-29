package repositories

import (
	"bonchDvach/pkg/handlers"
	"bonchDvach/pkg/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type BoardStorage struct {
	db *pgxpool.Pool
}

func NewBoardRepository(db *pgxpool.Pool) handlers.BoardRepository {
	return &BoardStorage{
		db: db,
	}
}

func (r *BoardStorage) GetBoards(ctx context.Context) ([]models.Board, error) {
	query := "SELECT id, name, description FROM boards"
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("cannot get boards: %w", err)
	}
	defer rows.Close()

	var boards []models.Board
	for rows.Next() {
		var b models.Board
		if err := rows.Scan(&b.ID, &b.Name, &b.Description); err != nil {
			return nil, fmt.Errorf("cannot scan rows: %w", err)
		}

		boards = append(boards, b)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while scanning rows: %w", err)
	}

	return boards, nil

}

func (r *BoardStorage) CreateBoard(ctx context.Context, name string, description string) error {
	query := "INSERT INTO boards (name, description) VALUES ($1, $2)"

	_, err := r.db.Exec(ctx, query, name, description)
	if err != nil {
		return err
	}

	return nil
}
