package usecase

import (
	"context"
	"github.com/SlavaShagalov/my-trello-backend/internal/auth"
	"github.com/SlavaShagalov/my-trello-backend/internal/models"
	pkgErrors "github.com/SlavaShagalov/my-trello-backend/internal/pkg/errors"
	pkgHasher "github.com/SlavaShagalov/my-trello-backend/internal/pkg/hasher"
	"github.com/SlavaShagalov/my-trello-backend/internal/sessions"
	"github.com/SlavaShagalov/my-trello-backend/internal/users"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"time"
)

type usecase struct {
	usersRepo    users.Repository
	sessionsRepo sessions.Repository
	hasher       pkgHasher.Hasher
	log          *zap.Logger
	tracer       trace.Tracer
}

func New(usersRepo users.Repository, sessionsRepo sessions.Repository, hasher pkgHasher.Hasher, log *zap.Logger, tracer trace.Tracer) auth.Usecase {
	return &usecase{
		usersRepo:    usersRepo,
		sessionsRepo: sessionsRepo,
		hasher:       hasher,
		log:          log,
		tracer:       tracer,
	}
}

func (uc *usecase) SignIn(ctx context.Context, params *auth.SignInParams) (models.User, string, error) {
	uc.log.Debug("HERE")
	ctx, span := uc.tracer.Start(ctx, "Usecase SignIn")
	time.Sleep(3 * time.Millisecond)
	defer span.End()

	user, err := uc.usersRepo.GetByUsername(ctx, params.Username)
	if err != nil {
		return models.User{}, "", err
	}
	time.Sleep(1 * time.Millisecond)

	if err = uc.hasher.CompareHashAndPassword(user.Password, params.Password); err != nil {
		return models.User{}, "", errors.Wrap(pkgErrors.ErrWrongLoginOrPassword, err.Error())
	}

	authToken, err := uc.sessionsRepo.Create(user.ID)
	if err != nil {
		return models.User{}, "", err
	}

	uc.log.Debug("Sign In", zap.Int("user_id", user.ID))
	return user, authToken, nil
}

func (uc *usecase) SignUp(params *auth.SignUpParams) (models.User, string, error) {
	_, err := uc.usersRepo.GetByUsername(context.TODO(), params.Username)
	if !errors.Is(err, pkgErrors.ErrUserNotFound) {
		if err != nil {
			return models.User{}, "", err
		}
		return models.User{}, "", pkgErrors.ErrUserAlreadyExists
	}

	hashedPassword, err := uc.hasher.GetHashedPassword(params.Password)
	if err != nil {
		return models.User{}, "", errors.Wrap(pkgErrors.ErrGetHashedPassword, err.Error())
	}

	repParams := &users.CreateParams{
		Name:           params.Name,
		Username:       params.Username,
		Email:          params.Email,
		HashedPassword: hashedPassword,
	}
	user, err := uc.usersRepo.Create(repParams)
	if err != nil {
		return models.User{}, "", err
	}

	authToken, err := uc.sessionsRepo.Create(user.ID)
	if err != nil {
		return models.User{}, "", err
	}

	uc.log.Debug("Sign Up", zap.Int("user_id", user.ID))
	return user, authToken, nil
}

func (uc *usecase) CheckAuth(userID int, authToken string) (int, error) {
	return uc.sessionsRepo.Get(userID, authToken)
}

func (uc *usecase) Logout(userID int, authToken string) error {
	err := uc.sessionsRepo.Delete(userID, authToken)
	if err != nil {
		return err
	}
	uc.log.Debug("Logout", zap.Int("user_id", userID))
	return nil
}
