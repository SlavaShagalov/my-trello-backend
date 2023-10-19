package workspaces

import (
	"github.com/SlavaShagalov/my-trello-backend/internal/models"
)

type CreateParams struct {
	Title       string
	Description string
	UserID      int
}

type FullUpdateParams struct {
	ID          int
	Title       string
	Description string
}

type PartialUpdateParams struct {
	ID                int
	Title             string
	UpdateTitle       bool
	Description       string
	UpdateDescription bool
}

type Repository interface {
	Create(params *CreateParams) (models.Workspace, error)
	List(userID int) ([]models.Workspace, error)
	Get(id int) (models.Workspace, error)
	FullUpdate(params *FullUpdateParams) (models.Workspace, error)
	PartialUpdate(params *PartialUpdateParams) (models.Workspace, error)
	Delete(id int) error
}
