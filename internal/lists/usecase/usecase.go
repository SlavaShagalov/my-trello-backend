package usecase

import (
	"github.com/SlavaShagalov/my-trello-backend/internal/lists"
	"github.com/SlavaShagalov/my-trello-backend/internal/models"
)

type usecase struct {
	repo lists.Repository
}

func NewUsecase(repo lists.Repository) lists.Usecase {
	return &usecase{repo: repo}
}

func (uc *usecase) Create(params *lists.CreateParams) (models.List, error) {
	return uc.repo.Create(params)
}

func (uc *usecase) List(userID int) ([]models.List, error) {
	return uc.repo.List(userID)
}

func (uc *usecase) Get(id int) (models.List, error) {
	return uc.repo.Get(id)
}

func (uc *usecase) FullUpdate(params *lists.FullUpdateParams) (models.List, error) {
	return uc.repo.FullUpdate(params)
}

func (uc *usecase) PartialUpdate(params *lists.PartialUpdateParams) (models.List, error) {
	return uc.repo.PartialUpdate(params)
}

func (uc *usecase) Delete(id int) error {
	return uc.repo.Delete(id)
}
