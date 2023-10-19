package usecase

import (
	"github.com/SlavaShagalov/my-trello-backend/internal/boards"
	"github.com/SlavaShagalov/my-trello-backend/internal/models"
)

type usecase struct {
	repo boards.Repository
}

func NewUsecase(repo boards.Repository) boards.Usecase {
	return &usecase{repo: repo}
}

func (uc *usecase) Create(params *boards.CreateParams) (models.Board, error) {
	return uc.repo.Create(params)
}

func (uc *usecase) List(userID int) ([]models.Board, error) {
	return uc.repo.List(userID)
}

func (uc *usecase) Get(id int) (models.Board, error) {
	return uc.repo.Get(id)
}

func (uc *usecase) FullUpdate(params *boards.FullUpdateParams) (models.Board, error) {
	return uc.repo.FullUpdate(params)
}

func (uc *usecase) PartialUpdate(params *boards.PartialUpdateParams) (models.Board, error) {
	return uc.repo.PartialUpdate(params)
}

func (uc *usecase) Delete(id int) error {
	return uc.repo.Delete(id)
}
