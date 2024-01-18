package usecase

import (
	"context"
	"github.com/SlavaShagalov/my-trello-backend/internal/auth"
	"github.com/SlavaShagalov/my-trello-backend/internal/models"
	pkgErrors "github.com/SlavaShagalov/my-trello-backend/internal/pkg/errors"
	pkgHasher "github.com/SlavaShagalov/my-trello-backend/internal/pkg/hasher"
	"github.com/SlavaShagalov/my-trello-backend/internal/pkg/opentel"
	"github.com/SlavaShagalov/my-trello-backend/internal/sessions"
	"github.com/SlavaShagalov/my-trello-backend/internal/users"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	componentName = "Auth Usecase"
)

type usecase struct {
	usersRepo    users.Repository
	sessionsRepo sessions.Repository
	hasher       pkgHasher.Hasher
	log          *zap.Logger
}

func New(usersRepo users.Repository, sessionsRepo sessions.Repository, hasher pkgHasher.Hasher, log *zap.Logger) auth.Usecase {
	return &usecase{
		usersRepo:    usersRepo,
		sessionsRepo: sessionsRepo,
		hasher:       hasher,
		log:          log,
	}
}

func (uc *usecase) SignIn(ctx context.Context, params *auth.SignInParams) (models.User, string, error) {
	ctx, span := opentel.Tracer.Start(ctx, componentName+" "+"SignIn")
	defer span.End()

	user, err := uc.usersRepo.GetByUsername(ctx, params.Username)
	if err != nil {
		return models.User{}, "", err
	}

	if err = uc.hasher.CompareHashAndPassword(ctx, user.Password, params.Password); err != nil {
		return models.User{}, "", errors.Wrap(pkgErrors.ErrWrongLoginOrPassword, err.Error())
	}

	authToken, err := uc.sessionsRepo.Create(ctx, user.ID)
	if err != nil {
		return models.User{}, "", err
	}

	uc.log.Debug("Sign In", zap.Int("user_id", user.ID))
	return user, authToken, nil
}

func (uc *usecase) SignUp(ctx context.Context, params *auth.SignUpParams) (models.User, string, error) {
	ctx, span := opentel.Tracer.Start(ctx, componentName+" "+"SignUp")
	defer span.End()

	_, err := uc.usersRepo.GetByUsername(ctx, params.Username)
	if !errors.Is(err, pkgErrors.ErrUserNotFound) {
		if err != nil {
			return models.User{}, "", err
		}
		return models.User{}, "", pkgErrors.ErrUserAlreadyExists
	}

	hashedPassword, err := uc.hasher.GetHashedPassword(ctx, params.Password)
	if err != nil {
		return models.User{}, "", errors.Wrap(pkgErrors.ErrGetHashedPassword, err.Error())
	}

	repParams := &users.CreateParams{
		Name:           params.Name,
		Username:       params.Username,
		Email:          params.Email,
		HashedPassword: hashedPassword,
	}
	user, err := uc.usersRepo.Create(ctx, repParams)
	if err != nil {
		return models.User{}, "", err
	}

	authToken, err := uc.sessionsRepo.Create(ctx, user.ID)
	if err != nil {
		return models.User{}, "", err
	}

	uc.log.Debug("Sign Up", zap.Int("user_id", user.ID))
	return user, authToken, nil
}

func (uc *usecase) CheckAuth(ctx context.Context, userID int, authToken string) (int, error) {
	ctx, span := opentel.Tracer.Start(ctx, componentName+" "+"CheckAuth")
	defer span.End()

	return uc.sessionsRepo.Get(ctx, userID, authToken)
}

func (uc *usecase) Logout(ctx context.Context, userID int, authToken string) error {
	ctx, span := opentel.Tracer.Start(ctx, componentName+" "+"Logout")
	defer span.End()

	err := uc.sessionsRepo.Delete(ctx, userID, authToken)
	if err != nil {
		return err
	}
	uc.log.Debug("Logout", zap.Int("user_id", userID))
	return nil
}
