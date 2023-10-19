package models

import "time"

type List struct {
	ID        int       `json:"id"`
	BoardID   int       `json:"board_id"`
	Title     string    `json:"title"`
	Position  int       `json:"position"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
