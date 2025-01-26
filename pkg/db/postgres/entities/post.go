package entities

import (
	"bonchDvach/pkg/handlers"
	"bonchDvach/pkg/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PostStorage struct {
	db *pgxpool.Pool
}

func NewPostRepository(db *pgxpool.Pool) handlers.PostRepository {
	return PostStorage{
		db: db,
	}
}

func (r PostStorage) GetAllPosts(ctx context.Context, threadID string) ([]models.Post, error) {
	query := "SELECT * FROM posts WHERE thread_id = $1"

	rows, err := r.db.Query(ctx, query, threadID)
	if err != nil {
		return nil, fmt.Errorf("cannot get posts: %w", err)
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(&p.ID, &p.ThreadID, &p.Content); err != nil {
			return nil, fmt.Errorf("cannot scan rows: %w", err)
		}
		posts = append(posts, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while scanning rows: %w", err)
	}

	return []models.Post{}, nil
}

func (r PostStorage) CreatePost(ctx context.Context, threadID string, content string) error {
	query := "INSERT INTO posts (thread_id, content) VALUES ($1, $2)"

	if _, err := r.db.Exec(ctx, query, threadID, content); err != nil {
		return err
	}

	return nil
}
