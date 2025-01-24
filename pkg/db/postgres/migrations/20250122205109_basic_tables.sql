-- +goose Up
-- +goose StatementBegin
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
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS boards CASCADE;
DROP TABLE IF EXISTS threads CASCADE;
DROP TABLE IF EXISTS posts CASCADE;
-- +goose StatementEnd
