package http

import (
	"github.com/SlavaShagalov/my-trello-backend/internal/models"
	"time"
)

//go:generate easyjson -all -snake_case models.go

// API requests
type createRequest struct {
	Title string `json:"title"`
}

type partialUpdateRequest struct {
	Title    *string `json:"title"`
	Position *int    `json:"position"`
	BoardID  *int    `json:"board_id"`
}

// API responses
type listResponse struct {
	Lists []models.List `json:"lists"`
}

func newListResponse(lists []models.List) *listResponse {
	return &listResponse{
		Lists: lists,
	}
}

type createResponse struct {
	ID        int       `json:"id"`
	BoardID   int       `json:"board_id"`
	Title     string    `json:"title"`
	Position  int       `json:"position"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func newCreateResponse(list *models.List) *createResponse {
	return &createResponse{
		ID:        list.ID,
		BoardID:   list.BoardID,
		Title:     list.Title,
		Position:  list.Position,
		CreatedAt: list.CreatedAt,
		UpdatedAt: list.UpdatedAt,
	}
}

type getResponse struct {
	ID        int       `json:"id"`
	BoardID   int       `json:"board_id"`
	Title     string    `json:"title"`
	Position  int       `json:"position"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func newGetResponse(list *models.List) *getResponse {
	return &getResponse{
		ID:        list.ID,
		BoardID:   list.BoardID,
		Title:     list.Title,
		Position:  list.Position,
		CreatedAt: list.CreatedAt,
		UpdatedAt: list.UpdatedAt,
	}
}
