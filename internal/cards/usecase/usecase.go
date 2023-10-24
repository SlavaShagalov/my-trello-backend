package usecase

import (
	"github.com/SlavaShagalov/my-trello-backend/internal/cards"
	"github.com/SlavaShagalov/my-trello-backend/internal/models"
)

type usecase struct {
	repo cards.Repository
}

func NewUsecase(repo cards.Repository) cards.Usecase {
	return &usecase{repo: repo}
}

func (uc *usecase) Create(params *cards.CreateParams) (models.Card, error) {
	return uc.repo.Create(params)
}

func (uc *usecase) ListByList(userID int) ([]models.Card, error) {
	return uc.repo.ListByList(userID)
}

func (uc *usecase) ListByTitle(title string, userID int) ([]models.Card, error) {
	return uc.repo.ListByTitle(title, userID)
}

func (uc *usecase) Get(id int) (models.Card, error) {
	return uc.repo.Get(id)
}

func (uc *usecase) FullUpdate(params *cards.FullUpdateParams) (models.Card, error) {
	return uc.repo.FullUpdate(params)
}

func (uc *usecase) PartialUpdate(params *cards.PartialUpdateParams) (models.Card, error) {
	return uc.repo.PartialUpdate(params)
}

func (uc *usecase) Delete(id int) error {
	return uc.repo.Delete(id)
}
