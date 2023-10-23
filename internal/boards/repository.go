package boards

import "github.com/SlavaShagalov/my-trello-backend/internal/models"

type CreateParams struct {
	Title       string
	Description string
	WorkspaceID int
}

type FullUpdateParams struct {
	ID          int
	Title       string
	Description string
	WorkspaceID int
}

type PartialUpdateParams struct {
	ID                int
	Title             string
	UpdateTitle       bool
	Description       string
	UpdateDescription bool
	WorkspaceID       int
	UpdateWorkspaceID bool
}

type Repository interface {
	Create(params *CreateParams) (models.Board, error)
	List(workspaceID int) ([]models.Board, error)
	ListByTitle(title string, userID int) ([]models.Board, error)
	Get(id int) (models.Board, error)
	FullUpdate(params *FullUpdateParams) (models.Board, error)
	PartialUpdate(params *PartialUpdateParams) (models.Board, error)
	UpdateBackground(id int, background string) error
	Delete(id int) error
}
