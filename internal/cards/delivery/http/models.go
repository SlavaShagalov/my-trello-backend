package http

import (
	"github.com/SlavaShagalov/my-trello-backend/internal/models"
	"time"
)

//go:generate easyjson -all -snake_case models.go

// API requests
type CreateRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type PartialUpdateRequest struct {
	Title    *string `json:"title"`
	Content  *string `json:"content"`
	Position *int    `json:"position"`
	ListID   *int    `json:"list_id"`
}

// API responses
type CardResponse struct {
	Cards []models.Card `json:"cards"`
}

func newListResponse(cards []models.Card) *CardResponse {
	return &CardResponse{
		Cards: cards,
	}
}

type CreateResponse struct {
	ID        int       `json:"id"`
	ListID    int       `json:"list_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Position  int       `json:"position"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func newCreateResponse(card *models.Card) *CreateResponse {
	return &CreateResponse{
		ID:        card.ID,
		ListID:    card.ListID,
		Title:     card.Title,
		Content:   card.Content,
		Position:  card.Position,
		CreatedAt: card.CreatedAt,
		UpdatedAt: card.UpdatedAt,
	}
}

type getResponse struct {
	ID        int       `json:"id"`
	ListID    int       `json:"list_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Position  int       `json:"position"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func newGetResponse(card *models.Card) *getResponse {
	return &getResponse{
		ID:        card.ID,
		ListID:    card.ListID,
		Title:     card.Title,
		Content:   card.Content,
		Position:  card.Position,
		CreatedAt: card.CreatedAt,
		UpdatedAt: card.UpdatedAt,
	}
}
