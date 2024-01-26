package usecase

import (
	"context"
	"github.com/SlavaShagalov/my-trello-backend/internal/boards"
	"github.com/SlavaShagalov/my-trello-backend/internal/images"
	"github.com/SlavaShagalov/my-trello-backend/internal/models"
	"github.com/SlavaShagalov/my-trello-backend/internal/pkg/opentel"
	"github.com/google/uuid"
	"path/filepath"
)

const (
	backgroundsFolder = "backgrounds"

	componentName = "Boards Usecase"
)

type usecase struct {
	repo    boards.Repository
	imgRepo images.Repository
}

func New(repo boards.Repository, imgRepo images.Repository) boards.Usecase {
	return &usecase{
		repo:    repo,
		imgRepo: imgRepo,
	}
}

func (uc *usecase) Create(ctx context.Context, params *boards.CreateParams) (models.Board, error) {
	ctx, span := opentel.Tracer.Start(ctx, componentName+" "+"Create")
	defer span.End()

	return uc.repo.Create(ctx, params)
}

func (uc *usecase) ListByWorkspace(ctx context.Context, userID int) ([]models.Board, error) {
	ctx, span := opentel.Tracer.Start(ctx, componentName+" "+"ListByWorkspace")
	defer span.End()

	return uc.repo.List(ctx, userID)
}

func (uc *usecase) ListByTitle(ctx context.Context, title string, userID int) ([]models.Board, error) {
	ctx, span := opentel.Tracer.Start(ctx, componentName+" "+"ListByTitle")
	defer span.End()

	return uc.repo.ListByTitle(ctx, title, userID)
}

func (uc *usecase) Get(ctx context.Context, id int) (models.Board, error) {
	ctx, span := opentel.Tracer.Start(ctx, componentName+" "+"Get")
	defer span.End()

	return uc.repo.Get(ctx, id)
}

func (uc *usecase) FullUpdate(ctx context.Context, params *boards.FullUpdateParams) (models.Board, error) {
	ctx, span := opentel.Tracer.Start(ctx, componentName+" "+"FullUpdate")
	defer span.End()

	return uc.repo.FullUpdate(ctx, params)
}

func (uc *usecase) PartialUpdate(ctx context.Context, params *boards.PartialUpdateParams) (models.Board, error) {
	ctx, span := opentel.Tracer.Start(ctx, componentName+" "+"PartialUpdate")
	defer span.End()

	return uc.repo.PartialUpdate(ctx, params)
}

func (uc *usecase) UpdateBackground(ctx context.Context, id int, imgData []byte, filename string) (*models.Board, error) {
	ctx, span := opentel.Tracer.Start(ctx, componentName+" "+"UpdateBackground")
	defer span.End()

	board, err := uc.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if board.Background == nil {
		imgName := backgroundsFolder + "/" + uuid.NewString() + filepath.Ext(filename)
		imgPath, err := uc.imgRepo.Create(imgName, imgData)
		if err == nil {
			err = uc.repo.UpdateBackground(ctx, id, imgPath)
			if err == nil {
				board.Background = &imgPath
			}
		}
	} else {
		err = uc.imgRepo.Update(*board.Background, imgData)
	}

	return &board, err
}

func (uc *usecase) Delete(ctx context.Context, id int) error {
	ctx, span := opentel.Tracer.Start(ctx, componentName+" "+"Delete")
	defer span.End()

	return uc.repo.Delete(ctx, id)
}
