package models

import "time"

type Board struct {
	ID          int       `json:"id"`
	WorkspaceID int       `json:"workspace_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Background  *string   `json:"background"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
