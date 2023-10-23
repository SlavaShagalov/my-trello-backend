package http

import (
	"github.com/SlavaShagalov/my-trello-backend/internal/models"
	"time"
)

//go:generate easyjson -all -snake_case models.go

// API requests
type createRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type partialUpdateRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

// API responses
type listResponse struct {
	Boards []models.Board `json:"boards"`
}

func newListResponse(boards []models.Board) *listResponse {
	return &listResponse{
		Boards: boards,
	}
}

type createResponse struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func newCreateResponse(board *models.Board) *createResponse {
	return &createResponse{
		ID:          board.ID,
		Title:       board.Title,
		Description: board.Description,
		CreatedAt:   board.CreatedAt,
		UpdatedAt:   board.UpdatedAt,
	}
}

type getResponse struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Background  *string   `json:"background"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func newGetResponse(board *models.Board) *getResponse {
	return &getResponse{
		ID:          board.ID,
		Title:       board.Title,
		Description: board.Description,
		Background:  board.Background,
		CreatedAt:   board.CreatedAt,
		UpdatedAt:   board.UpdatedAt,
	}
}
