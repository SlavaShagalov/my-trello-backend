package models

import "time"

type Card struct {
	ID        int       `json:"id"`
	ListID    int       `json:"list_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Position  int       `json:"position"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
