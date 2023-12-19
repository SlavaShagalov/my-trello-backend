package integration

import (
	"context"
	"database/sql"
	"github.com/SlavaShagalov/my-trello-backend/internal/models"
	"github.com/SlavaShagalov/my-trello-backend/internal/users"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/SlavaShagalov/my-trello-backend/internal/pkg/config"

	pkgErrors "github.com/SlavaShagalov/my-trello-backend/internal/pkg/errors"
	pkgHasher "github.com/SlavaShagalov/my-trello-backend/internal/pkg/hasher/bcrypt"
	pkgZap "github.com/SlavaShagalov/my-trello-backend/internal/pkg/log/zap"
	pkgDb "github.com/SlavaShagalov/my-trello-backend/internal/pkg/storages"

	pkgAuth "github.com/SlavaShagalov/my-trello-backend/internal/auth"
	authUC "github.com/SlavaShagalov/my-trello-backend/internal/auth/usecase"
	sessionsRepository "github.com/SlavaShagalov/my-trello-backend/internal/sessions/repository/redis"
	usersRepository "github.com/SlavaShagalov/my-trello-backend/internal/users/repository/postgres"
)

type AuthSuite struct {
	suite.Suite
	db        *sql.DB
	rdb       *redis.Client
	log       *zap.Logger
	usersRepo users.Repository
	uc        pkgAuth.Usecase
}

func (s *AuthSuite) SetupSuite() {
	s.log = pkgZap.NewTestLogger()

	var err error
	config.SetTestPostgresConfig()
	s.db, err = pkgDb.NewPostgres(s.log)
	s.Require().NoError(err)

	config.SetTestRedisConfig()
	ctx := context.Background()
	s.rdb, err = pkgDb.NewRedis(s.log, ctx)
	s.Require().NoError(err)

	s.usersRepo = usersRepository.New(s.db, s.log)
	sessionsRepo := sessionsRepository.New(s.rdb, ctx, s.log)
	hasher := pkgHasher.New()
	s.uc = authUC.New(s.usersRepo, sessionsRepo, hasher, s.log)
}

func (s *AuthSuite) TearDownSuite() {
	err := s.db.Close()
	s.Require().NoError(err)
	s.log.Info("DB connection closed")

	err = s.rdb.Close()
	s.Require().NoError(err)
	s.log.Info("Redis connection closed")

	_ = s.log.Sync()
}

func (s *AuthSuite) TestSignIn() {
	type testCase struct {
		params *pkgAuth.SignInParams
		user   models.User
		err    error
	}

	tests := map[string]testCase{
		"normal": {
			params: &pkgAuth.SignInParams{
				Username: "slava",
				Password: "12345678",
			},
			user: models.User{
				ID:       1,
				Username: "slava",
				Password: "$2a$10$A4Ab/cuy/oLNvm4VxGoO/ezKL.fiew5e.eKTkUOWIVxoBh8XFO4lS",
				Email:    "slava@vk.com",
				Name:     "Slava",
			},
			err: nil,
		},
		"wrong password": {
			params: &pkgAuth.SignInParams{
				Username: "slava",
				Password: "12345679",
			},
			user: models.User{},
			err:  pkgErrors.ErrWrongLoginOrPassword,
		},
		"user not found": {
			params: &pkgAuth.SignInParams{
				Username: "noname",
				Password: "12345678",
			},
			user: models.User{},
			err:  pkgErrors.ErrUserNotFound,
		},
	}

	for name, test := range tests {
		s.Run(name, func() {
			user, authToken, err := s.uc.SignIn(test.params)
			assert.ErrorIs(s.T(), err, test.err, "unexpected error")

			assert.Equal(s.T(), test.user.ID, user.ID, "incorrect user ID")
			assert.Equal(s.T(), test.user.Username, user.Username, "incorrect Username")
			assert.Equal(s.T(), test.user.Password, user.Password, "incorrect Password")
			assert.Equal(s.T(), test.user.Email, user.Email, "incorrect Email")
			assert.Equal(s.T(), test.user.Name, user.Name, "incorrect Name")

			if err == nil {
				assert.NotEmpty(s.T(), authToken, "incorrect AuthToken")

				_, err = s.uc.CheckAuth(user.ID, authToken)
				assert.NoError(s.T(), err, "unexpected unauthorized")

				err = s.uc.Logout(user.ID, authToken)
				assert.NoError(s.T(), err, "failed to logout user")
			}
		})
	}
}

func (s *AuthSuite) TestSignUp() {
	type testCase struct {
		params *pkgAuth.SignUpParams
		user   models.User
		err    error
	}

	tests := map[string]testCase{
		"normal": {
			params: &pkgAuth.SignUpParams{
				Name:     "New User",
				Username: "new_user",
				Email:    "new_user@vk.com",
				Password: "12345678",
			},
			user: models.User{
				Username: "new_user",
				Email:    "new_user@vk.com",
				Name:     "New User",
			},
			err: nil,
		},
		"user with such username already exists": {
			params: &pkgAuth.SignUpParams{
				Name:     "New Slava",
				Username: "slava",
				Email:    "new_slava@vk.com",
				Password: "123456789",
			},
			user: models.User{},
			err:  pkgErrors.ErrUserAlreadyExists,
		},
	}

	for name, test := range tests {
		s.Run(name, func() {
			user, authToken, err := s.uc.SignUp(test.params)
			assert.ErrorIs(s.T(), err, test.err, "unexpected error")

			assert.Equal(s.T(), test.user.Username, user.Username, "incorrect Username")
			assert.Equal(s.T(), test.user.Email, user.Email, "incorrect Email")
			assert.Equal(s.T(), test.user.Name, user.Name, "incorrect Name")

			if err == nil {
				assert.NotEmpty(s.T(), authToken, "incorrect AuthToken")

				_, err = s.uc.CheckAuth(user.ID, authToken)
				assert.NoError(s.T(), err, "unexpected unauthorized")

				err = s.uc.Logout(user.ID, authToken)
				assert.NoError(s.T(), err, "failed to logout user")

				err = s.usersRepo.Delete(user.ID)
				assert.NoError(s.T(), err, "failed to delete user")
			}
		})
	}
}

func (s *AuthSuite) TestCheckAuth() {
	type testCase struct {
		userID    int
		authToken string
		err       error
	}

	// prepare session for tests
	user, validAuthToken, err := s.uc.SignIn(&pkgAuth.SignInParams{
		Username: "slava",
		Password: "12345678",
	})
	assert.NoError(s.T(), err, "unexpected error")

	tests := map[string]testCase{
		"normal": {
			userID:    1,
			authToken: validAuthToken,
			err:       nil,
		},
		"session not found due to incorrect token": {
			userID:    1,
			authToken: "invalid_token",
			err:       pkgErrors.ErrSessionNotFound,
		},
		"session not found due to incorrect user id": {
			userID:    2,
			authToken: validAuthToken,
			err:       pkgErrors.ErrSessionNotFound,
		},
	}

	for name, test := range tests {
		s.Run(name, func() {
			userID, err := s.uc.CheckAuth(test.userID, test.authToken)
			assert.ErrorIs(s.T(), err, test.err, "unexpected error")

			if err == nil {
				assert.Equal(s.T(), test.userID, userID, "incorrect user ID")
			}
		})
	}

	// delete prepared session
	err = s.uc.Logout(user.ID, validAuthToken)
	assert.NoError(s.T(), err, "failed to logout user")
}

func (s *AuthSuite) TestLogout() {
	type testCase struct {
		userID    int
		authToken string
		err       error
	}

	// prepare session for tests
	user, validAuthToken, err := s.uc.SignIn(&pkgAuth.SignInParams{
		Username: "slava",
		Password: "12345678",
	})
	assert.NoError(s.T(), err, "unexpected error")

	tests := map[string]testCase{
		"session not found due to incorrect token": {
			userID:    1,
			authToken: "invalid_token",
			err:       pkgErrors.ErrSessionNotFound,
		},
		"session not found due to incorrect user id": {
			userID:    2,
			authToken: validAuthToken,
			err:       pkgErrors.ErrSessionNotFound,
		},
		"normal": {
			userID:    1,
			authToken: validAuthToken,
			err:       nil,
		},
	}

	for name, test := range tests {
		s.Run(name, func() {
			err = s.uc.Logout(test.userID, test.authToken)
			assert.ErrorIs(s.T(), err, test.err, "unexpected error")

			if err == nil {
				_, err = s.uc.CheckAuth(user.ID, test.authToken)
				assert.ErrorIs(s.T(), err, pkgErrors.ErrSessionNotFound, "unexpected error")
			}
		})
	}
}

func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(AuthSuite))
}
