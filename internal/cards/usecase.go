package cards

import (
	"github.com/SlavaShagalov/my-trello-backend/internal/models"
)

type Usecase interface {
	Create(params *CreateParams) (models.Card, error)
	List(listID int) ([]models.Card, error)
	Get(id int) (models.Card, error)
	FullUpdate(params *FullUpdateParams) (models.Card, error)
	PartialUpdate(params *PartialUpdateParams) (models.Card, error)
	Delete(id int) error
}
