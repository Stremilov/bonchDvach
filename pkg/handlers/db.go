package handlers

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

var db *sql.DB

func InitDB() {
	var err error

	connStr := fmt.Sprintf("user=levstremilov password=postgres dbname=bonchdvach sslmode=disable")

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping database: %v", err)
	}

	query := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			ip VARCHAR(255) NOT NULL
		);
		CREATE TABLE IF NOT EXISTS boards (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS threads (
			id SERIAL PRIMARY KEY,
			board_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			FOREIGN KEY (board_id) REFERENCES boards (id)
		);
		CREATE TABLE IF NOT EXISTS posts (
			id SERIAL PRIMARY KEY,
			thread_id INTEGER NOT NULL,
			content TEXT NOT NULL,
			FOREIGN KEY (thread_id) REFERENCES threads (id)
		);
	`

	_, err = db.Exec(query)
	if err != nil {
		log.Fatal("Failed to create tables: %v", err)
	}

	fmt.Println("Successfully connected to database and tables created")
}
