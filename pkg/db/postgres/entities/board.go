package entities

import (
	"bonchDvach/pkg/models"
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

type BoardRepository struct {
	db *pgxpool.Pool
}

func NewBoardRepository(db *pgxpool.Pool) *BoardRepository {
	return &BoardRepository{
		db: db,
	}
}

func (r *BoardRepository) GetBoards(ctx context.Context) ([]models.Board, error) {
	query := "SELECT id, name, description FROM boards"
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var boards []models.Board
	for rows.Next() {
		var b models.Board
		if err := rows.Scan(&b.ID, &b.Name, &b.Description); err != nil {
			log.Println(err)
		}

		boards = append(boards, b)
	}

	return boards, nil

}

func (r *BoardRepository) CreateBoard(ctx context.Context, name string, description string) error {
	query := "INSERT INTO boards (name, description) VALUES ($1, $2)"

	_, err := r.db.Exec(ctx, query, name, description)
	if err != nil {
		return err
	}

	return nil
}
