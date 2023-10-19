package users

import (
	"github.com/SlavaShagalov/my-trello-backend/internal/models"
)

type Usecase interface {
	List() ([]models.User, error)
	Get(id int) (models.User, error)
	GetByUsername(username string) (models.User, error)
	FullUpdate(params *FullUpdateParams) (models.User, error)
	PartialUpdate(params *PartialUpdateParams) (models.User, error)
	Delete(id int) error
}
