package cards

import "github.com/SlavaShagalov/my-trello-backend/internal/models"

type CreateParams struct {
	Title   string
	Content string
	ListID  int
}

type FullUpdateParams struct {
	ID       int
	Title    string
	Content  string
	Position int
	ListID   int
}

type PartialUpdateParams struct {
	ID             int
	Title          string
	UpdateTitle    bool
	Content        string
	UpdateContent  bool
	Position       int
	UpdatePosition bool
	ListID         int
	UpdateListID   bool
}

type Repository interface {
	Create(params *CreateParams) (models.Card, error)
	ListByList(listID int) ([]models.Card, error)
	ListByTitle(title string, userID int) ([]models.Card, error)
	Get(id int) (models.Card, error)
	FullUpdate(params *FullUpdateParams) (models.Card, error)
	PartialUpdate(params *PartialUpdateParams) (models.Card, error)
	Delete(id int) error
}
