package boards

import (
	"context"
	"github.com/SlavaShagalov/my-trello-backend/internal/models"
)

type Usecase interface {
	Create(ctx context.Context, params *CreateParams) (models.Board, error)
	ListByWorkspace(ctx context.Context, workspaceID int) ([]models.Board, error)
	ListByTitle(ctx context.Context, title string, userID int) ([]models.Board, error)
	Get(ctx context.Context, id int) (models.Board, error)
	FullUpdate(ctx context.Context, params *FullUpdateParams) (models.Board, error)
	PartialUpdate(ctx context.Context, params *PartialUpdateParams) (models.Board, error)
	UpdateBackground(ctx context.Context, id int, imgData []byte, filename string) (*models.Board, error)
	Delete(ctx context.Context, id int) error
}
