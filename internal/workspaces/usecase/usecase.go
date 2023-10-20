package usecase

import (
	"github.com/SlavaShagalov/my-trello-backend/internal/models"
	"github.com/SlavaShagalov/my-trello-backend/internal/pkg/constants"
	pkgErrors "github.com/SlavaShagalov/my-trello-backend/internal/pkg/errors"
	"github.com/SlavaShagalov/my-trello-backend/internal/workspaces"
)

type usecase struct {
	rep workspaces.Repository
}

func New(rep workspaces.Repository) workspaces.Usecase {
	return &usecase{rep: rep}
}

func (uc *usecase) Create(params *workspaces.CreateParams) (models.Workspace, error) {
	if err := validateTitle(params.Title); err != nil {
		return models.Workspace{}, err
	} else if err = validateDescription(params.Description); err != nil {
		return models.Workspace{}, err
	}

	return uc.rep.Create(params)
}

func (uc *usecase) List(userID int) ([]models.Workspace, error) {
	return uc.rep.List(userID)
}

func (uc *usecase) Get(id int) (models.Workspace, error) {
	return uc.rep.Get(id)
}

func (uc *usecase) FullUpdate(params *workspaces.FullUpdateParams) (models.Workspace, error) {
	if err := validateTitle(params.Title); err != nil {
		return models.Workspace{}, err
	} else if err = validateDescription(params.Description); err != nil {
		return models.Workspace{}, err
	}

	return uc.rep.FullUpdate(params)
}

func (uc *usecase) PartialUpdate(params *workspaces.PartialUpdateParams) (models.Workspace, error) {
	if params.UpdateTitle {
		if err := validateTitle(params.Title); err != nil {
			return models.Workspace{}, err
		}
	} else if params.UpdateDescription {
		if err := validateDescription(params.Description); err != nil {
			return models.Workspace{}, err
		}
	}

	return uc.rep.PartialUpdate(params)
}

func (uc *usecase) Delete(id int) error {
	return uc.rep.Delete(id)
}

func validateTitle(title string) error {
	if len(title) > constants.MaxListTitleLen {
		return pkgErrors.ErrTooLongListTitle
	}
	return nil
}

func validateDescription(description string) error {
	if len(description) > constants.MaxListDescriptionLen {
		return pkgErrors.ErrTooLongListDescription
	}
	return nil
}
