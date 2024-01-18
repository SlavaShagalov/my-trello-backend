package users

import (
	"context"
	"github.com/SlavaShagalov/my-trello-backend/internal/models"
)

type CreateParams struct {
	Name           string
	Username       string
	Email          string
	HashedPassword string
}

type FullUpdateParams struct {
	ID       int
	Username string
	Email    string
	Name     string
}

type PartialUpdateParams struct {
	ID             int
	Username       string
	UpdateUsername bool
	Email          string
	UpdateEmail    bool
	Name           string
	UpdateName     bool
}

type Repository interface {
	Create(ctx context.Context, params *CreateParams) (models.User, error)

	List() ([]models.User, error)
	Get(id int) (models.User, error)
	GetByUsername(ctx context.Context, username string) (models.User, error)

	FullUpdate(params *FullUpdateParams) (models.User, error)
	PartialUpdate(params *PartialUpdateParams) (models.User, error)
	UpdateAvatar(id int, avatar string) error

	Delete(id int) error

	Exists(id int) (bool, error)
}
