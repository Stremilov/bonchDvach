package models

type Post struct {
	ID       int    `json:"id"`
	ThreadID int    `json:"threadID"`
	Content  string `json:"content"`
}
