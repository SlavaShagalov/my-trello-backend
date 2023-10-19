package workspaces

import "github.com/SlavaShagalov/my-trello-backend/internal/models"

type Usecase interface {
	Create(params *CreateParams) (models.Workspace, error)
	List(userID int) ([]models.Workspace, error)
	Get(id int) (models.Workspace, error)
	FullUpdate(params *FullUpdateParams) (models.Workspace, error)
	PartialUpdate(params *PartialUpdateParams) (models.Workspace, error)
	Delete(id int) error
}
