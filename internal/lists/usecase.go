package lists

import (
	"github.com/SlavaShagalov/my-trello-backend/internal/models"
)

type Usecase interface {
	Create(params *CreateParams) (models.List, error)
	ListByBoard(boardID int) ([]models.List, error)
	ListByTitle(title string, userID int) ([]models.List, error)
	Get(id int) (models.List, error)
	FullUpdate(params *FullUpdateParams) (models.List, error)
	PartialUpdate(params *PartialUpdateParams) (models.List, error)
	Delete(id int) error
}
