package models

type Thread struct {
	ID      int    `json:"id"`
	BoardID int    `json:"boardID"`
	Title   string `json:"title"`
}
