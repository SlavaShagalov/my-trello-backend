package lists

import "github.com/SlavaShagalov/my-trello-backend/internal/models"

type CreateParams struct {
	Title   string
	BoardID int
}

type FullUpdateParams struct {
	ID       int
	Title    string
	Position int
	BoardID  int
}

type PartialUpdateParams struct {
	ID             int
	Title          string
	UpdateTitle    bool
	Position       int
	UpdatePosition bool
	BoardID        int
	UpdateBoardID  bool
}

type Repository interface {
	Create(params *CreateParams) (models.List, error)
	ListByBoard(boardID int) ([]models.List, error)
	ListByTitle(title string, userID int) ([]models.List, error)
	Get(id int) (models.List, error)
	FullUpdate(params *FullUpdateParams) (models.List, error)
	PartialUpdate(params *PartialUpdateParams) (models.List, error)
	Delete(id int) error
}
