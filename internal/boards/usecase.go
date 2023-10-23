package boards

import (
	"github.com/SlavaShagalov/my-trello-backend/internal/models"
)

type Usecase interface {
	Create(params *CreateParams) (models.Board, error)
	ListByWorkspace(workspaceID int) ([]models.Board, error)
	ListByTitle(title string, userID int) ([]models.Board, error)
	Get(id int) (models.Board, error)
	FullUpdate(params *FullUpdateParams) (models.Board, error)
	PartialUpdate(params *PartialUpdateParams) (models.Board, error)
	UpdateBackground(id int, imgData []byte, filename string) (*models.Board, error)
	Delete(id int) error
}
