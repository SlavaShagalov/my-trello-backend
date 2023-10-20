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
	Workspaces []models.Workspace `json:"workspaces"`
}

func newListResponse(workspaces []models.Workspace) *listResponse {
	return &listResponse{
		Workspaces: workspaces,
	}
}

type createResponse struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func newCreateResponse(workspace *models.Workspace) *createResponse {
	return &createResponse{
		ID:          workspace.ID,
		Title:       workspace.Title,
		Description: workspace.Description,
		CreatedAt:   workspace.CreatedAt,
		UpdatedAt:   workspace.UpdatedAt,
	}
}

type getResponse struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func newGetResponse(workspace *models.Workspace) *getResponse {
	return &getResponse{
		ID:          workspace.ID,
		Title:       workspace.Title,
		Description: workspace.Description,
		CreatedAt:   workspace.CreatedAt,
		UpdatedAt:   workspace.UpdatedAt,
	}
}
