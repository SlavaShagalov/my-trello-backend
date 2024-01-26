package boards

import (
	"context"
	"github.com/SlavaShagalov/my-trello-backend/internal/models"
)

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
	Create(ctx context.Context, params *CreateParams) (models.Board, error)
	List(ctx context.Context, workspaceID int) ([]models.Board, error)
	ListByTitle(ctx context.Context, title string, userID int) ([]models.Board, error)
	Get(ctx context.Context, id int) (models.Board, error)
	FullUpdate(ctx context.Context, params *FullUpdateParams) (models.Board, error)
	PartialUpdate(ctx context.Context, params *PartialUpdateParams) (models.Board, error)
	UpdateBackground(ctx context.Context, id int, background string) error
	Delete(ctx context.Context, id int) error
}
