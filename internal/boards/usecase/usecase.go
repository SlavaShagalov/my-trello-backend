package usecase

import (
	"github.com/SlavaShagalov/my-trello-backend/internal/boards"
	"github.com/SlavaShagalov/my-trello-backend/internal/images"
	"github.com/SlavaShagalov/my-trello-backend/internal/models"
	"github.com/google/uuid"
	"path/filepath"
)

const (
	backgroundsFolder = "backgrounds"
)

type usecase struct {
	repo    boards.Repository
	imgRepo images.Repository
}

func NewUsecase(repo boards.Repository, imgRepo images.Repository) boards.Usecase {
	return &usecase{
		repo:    repo,
		imgRepo: imgRepo,
	}
}

func (uc *usecase) Create(params *boards.CreateParams) (models.Board, error) {
	return uc.repo.Create(params)
}

func (uc *usecase) ListByWorkspace(userID int) ([]models.Board, error) {
	return uc.repo.List(userID)
}

func (uc *usecase) ListByTitle(title string, userID int) ([]models.Board, error) {
	return uc.repo.ListByTitle(title, userID)
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

func (uc *usecase) UpdateBackground(id int, imgData []byte, filename string) (*models.Board, error) {
	board, err := uc.repo.Get(id)
	if err != nil {
		return nil, err
	}

	if board.Background == nil {
		imgName := backgroundsFolder + "/" + uuid.NewString() + filepath.Ext(filename)
		imgPath, err := uc.imgRepo.Create(imgName, imgData)
		if err == nil {
			err = uc.repo.UpdateBackground(id, imgPath)
			if err == nil {
				board.Background = &imgPath
			}
		}
	} else {
		err = uc.imgRepo.Update(*board.Background, imgData)
	}

	return &board, err
}

func (uc *usecase) Delete(id int) error {
	return uc.repo.Delete(id)
}
