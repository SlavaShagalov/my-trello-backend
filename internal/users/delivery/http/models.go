package http

import (
	"github.com/SlavaShagalov/my-trello-backend/internal/models"
	"time"
)

//go:generate easyjson -all -snake_case models.go

// API requests
type partialUpdateRequest struct {
	Username *string `json:"username"`
	Email    *string `json:"email"`
	Name     *string `json:"name"`
}

// API responses
type getResponse struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Avatar    *string   `json:"avatar"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func newGetResponse(user *models.User) *getResponse {
	return &getResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Name:      user.Name,
		Avatar:    user.Avatar,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
