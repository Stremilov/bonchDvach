package repositories

import (
	"bonchDvach/pkg/handlers"
	"bonchDvach/pkg/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type ThreadStorage struct {
	db *pgxpool.Pool
}

func NewThreadRepository(db *pgxpool.Pool) handlers.ThreadRepository {
	return ThreadStorage{
		db: db,
	}
}

// CreateThread implements db.ThreadRepository.
func (r ThreadStorage) CreateThread(ctx context.Context, title string, boardID string) error {
	query := "INSERT INTO threads (title, board_id) VALUES ($1, $2)"

	if _, err := r.db.Exec(ctx, query, title, boardID); err != nil {
		return err
	}

	return nil
}

// GetAllThreads implements db.ThreadRepository.
func (r ThreadStorage) GetAllThreads(ctx context.Context, boardID string) ([]models.Thread, error) {
	query := "SELECT id, board_id, title FROM threads WHERE board_id = $1"

	rows, err := r.db.Query(ctx, query, boardID)
	if err != nil {
		return nil, fmt.Errorf("cannot get threads: %w", err)
	}

	defer rows.Close()

	var threads []models.Thread
	for rows.Next() {
		var t models.Thread
		if err := rows.Scan(&t.ID, &t.BoardID, &t.Title); err != nil {
			return nil, fmt.Errorf("cannot scan rows: %w", err)
		}
		threads = append(threads, t)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error while scanning rows: %w", err)
	}

	return threads, nil

}
