package boards

import (
	"github.com/SlavaShagalov/my-trello-backend/internal/models"
)

type Usecase interface {
	Create(params *CreateParams) (models.Board, error)
	List(workspaceID int) ([]models.Board, error)
	Get(id int) (models.Board, error)
	FullUpdate(params *FullUpdateParams) (models.Board, error)
	PartialUpdate(params *PartialUpdateParams) (models.Board, error)
	Delete(id int) error
}
