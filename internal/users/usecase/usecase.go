package usecase

import (
	"github.com/SlavaShagalov/my-trello-backend/internal/models"
	"github.com/SlavaShagalov/my-trello-backend/internal/pkg/config"
	"github.com/SlavaShagalov/my-trello-backend/internal/pkg/constants"
	pkgErrors "github.com/SlavaShagalov/my-trello-backend/internal/pkg/errors"
	"github.com/SlavaShagalov/my-trello-backend/internal/users"
	"github.com/spf13/viper"
)

type usecase struct {
	rep users.Repository
}

func NewUsecase(rep users.Repository) users.Usecase {
	return &usecase{rep: rep}
}

func (uc *usecase) List() ([]models.User, error) {
	return uc.rep.List()
}

func (uc *usecase) Get(id int) (models.User, error) {
	return uc.rep.Get(id)
}

func (uc *usecase) GetByUsername(username string) (models.User, error) {
	return uc.rep.GetByUsername(username)
}

func (uc *usecase) FullUpdate(params *users.FullUpdateParams) (models.User, error) {
	if err := validateUsername(params.Username); err != nil {
		return models.User{}, err
	} else if err = validateName(params.Name); err != nil {
		return models.User{}, err
	}

	return uc.rep.FullUpdate(params)
}

func (uc *usecase) PartialUpdate(params *users.PartialUpdateParams) (models.User, error) {
	if params.UpdateUsername {
		if err := validateUsername(params.Username); err != nil {
			return models.User{}, err
		}
	} else if params.UpdateName {
		if err := validateName(params.Name); err != nil {
			return models.User{}, err
		}
	}

	return uc.rep.PartialUpdate(params)
}

func (uc *usecase) Delete(id int) error {
	return uc.rep.Delete(id)
}

func validateUsername(username string) error {
	if len(username) < viper.GetInt(config.MinUsernameLen) {
		return pkgErrors.ErrTooShortUsername
	} else if len(username) > viper.GetInt(config.MaxUsernameLen) {
		return pkgErrors.ErrTooLongUsername
	}
	return nil
}

func validateName(name string) error {
	if len(name) < constants.MinNameLen {
		return pkgErrors.ErrEmptyName
	} else if len(name) > constants.MaxNameLen {
		return pkgErrors.ErrTooLongName
	}
	return nil
}
