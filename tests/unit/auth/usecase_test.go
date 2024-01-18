package auth

import (
	"context"
	pkgZap "github.com/SlavaShagalov/my-trello-backend/internal/pkg/log/zap"
	"github.com/SlavaShagalov/my-trello-backend/internal/pkg/ot"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"

	pkgAuth "github.com/SlavaShagalov/my-trello-backend/internal/auth"
	authUsecase "github.com/SlavaShagalov/my-trello-backend/internal/auth/usecase"
	"github.com/SlavaShagalov/my-trello-backend/internal/models"
	pkgErrors "github.com/SlavaShagalov/my-trello-backend/internal/pkg/errors"
	"github.com/SlavaShagalov/my-trello-backend/internal/users"

	hasherMocks "github.com/SlavaShagalov/my-trello-backend/internal/pkg/hasher/mocks"
	sessionsMocks "github.com/SlavaShagalov/my-trello-backend/internal/sessions/mocks"
	usersMocks "github.com/SlavaShagalov/my-trello-backend/internal/users/mocks"
)

type AuthUsecaseSuite struct {
	suite.Suite
	logger *zap.Logger
	tp     *sdktrace.TracerProvider
	tracer trace.Tracer
	ctx    context.Context
}

func (s *AuthUsecaseSuite) BeforeAll(t provider.T) {
	t.WithNewStep("SetupSuite step", func(ctx provider.StepCtx) {})

	s.ctx = context.Background()

	s.logger = pkgZap.NewDevelopLogger()

	// Set up OpenTelemetry.
	exp, err := ot.NewOTLPExporter(s.ctx)
	if err != nil {
		s.logger.Error("Failed to initialize exporter", zap.Error(err))
	}

	s.tp = ot.NewTraceProvider(exp)
	otel.SetTracerProvider(s.tp)

	s.tracer = s.tp.Tracer("test")
	s.logger.Info("OpenTelemetry setup")
}

func (s *AuthUsecaseSuite) AfterAll(t provider.T) {
	t.WithNewStep("TearDownSuite step", func(ctx provider.StepCtx) {})

	_ = s.logger.Sync()

	_ = s.tp.Shutdown(s.ctx)
}

func (s *AuthUsecaseSuite) BeforeEach(t provider.T) {
	t.WithNewStep("SetupTest step", func(ctx provider.StepCtx) {})
}

func (s *AuthUsecaseSuite) AfterEach(t provider.T) {
	t.WithNewStep("TearDownTest step", func(ctx provider.StepCtx) {})
}

func (s *AuthUsecaseSuite) TestSignIn(t provider.T) {
	type fields struct {
		usersRepo    *usersMocks.MockRepository
		sessionsRepo *sessionsMocks.MockRepository
		hasher       *hasherMocks.MockHasher
		params       *pkgAuth.SignInParams
		user         *models.User
		authToken    string
	}

	type testCase struct {
		prepare   func(f *fields)
		params    *pkgAuth.SignInParams
		user      models.User
		authToken string
		err       error
	}

	tests := map[string]testCase{
		"normal": {
			prepare: func(f *fields) {
				gomock.InOrder(
					f.usersRepo.EXPECT().GetByUsername(gomock.Any(), f.params.Username).Return(*f.user, nil),
					f.hasher.EXPECT().CompareHashAndPassword(gomock.Any(), f.user.Password, f.params.Password).Return(nil),
					f.sessionsRepo.EXPECT().Create(gomock.Any(), f.user.ID).Return(f.authToken, nil),
				)
			},
			params: &pkgAuth.SignInParams{
				Username: "slava",
				Password: "1234",
			},
			user: models.User{
				ID:       21,
				Username: "slava",
				Password: "hash",
				Email:    "slava@vk.com",
				Name:     "Slava",
			},
			authToken: "auth_token",
			err:       nil,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t provider.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				usersRepo:    usersMocks.NewMockRepository(ctrl),
				sessionsRepo: sessionsMocks.NewMockRepository(ctrl),
				hasher:       hasherMocks.NewMockHasher(ctrl),
				params:       test.params,
				user:         &test.user,
				authToken:    test.authToken,
			}
			if test.prepare != nil {
				test.prepare(&f)
			}

			uc := authUsecase.New(f.usersRepo, f.sessionsRepo, f.hasher, s.logger, s.tracer)
			user, authToken, err := uc.SignIn(context.Background(), test.params)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if user != test.user {
				t.Errorf("\nExpected: %v\nGot: %v", test.user, user)
			}
			if authToken != test.authToken {
				t.Errorf("\nExpected: %v\nGot: %v", test.authToken, authToken)
			}
		})
	}
}

func (s *AuthUsecaseSuite) TestSignUp(t provider.T) {
	type fields struct {
		usersRepo    *usersMocks.MockRepository
		sessionsRepo *sessionsMocks.MockRepository
		hasher       *hasherMocks.MockHasher
		params       *pkgAuth.SignUpParams
		user         *models.User
		authToken    string
	}

	type testCase struct {
		prepare   func(f *fields)
		params    *pkgAuth.SignUpParams
		user      models.User
		authToken string
		err       error
	}

	tests := map[string]testCase{
		"normal": {
			prepare: func(f *fields) {
				gomock.InOrder(
					f.usersRepo.EXPECT().GetByUsername(gomock.Any(), f.params.Username).
						Return(models.User{}, pkgErrors.ErrUserNotFound),
					f.hasher.EXPECT().GetHashedPassword(gomock.Any(), f.params.Password).Return(f.user.Password, nil),
					f.usersRepo.EXPECT().Create(gomock.Any(), &users.CreateParams{
						Name:           f.params.Name,
						Username:       f.params.Username,
						Email:          f.params.Email,
						HashedPassword: f.user.Password,
					}).Return(*f.user, nil),
					f.sessionsRepo.EXPECT().Create(gomock.Any(), f.user.ID).Return(f.authToken, nil),
				)
			},
			params: &pkgAuth.SignUpParams{
				Name:     "Slava",
				Username: "slava",
				Email:    "slava@vk.com",
				Password: "1234",
			},
			user: models.User{
				ID:       21,
				Username: "slava",
				Password: "hash",
				Email:    "slava@vk.com",
				Name:     "Slava",
			},
			authToken: "auth_token",
			err:       nil,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t provider.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				usersRepo:    usersMocks.NewMockRepository(ctrl),
				sessionsRepo: sessionsMocks.NewMockRepository(ctrl),
				hasher:       hasherMocks.NewMockHasher(ctrl),
				params:       test.params,
				user:         &test.user,
				authToken:    test.authToken,
			}
			if test.prepare != nil {
				test.prepare(&f)
			}

			uc := authUsecase.New(f.usersRepo, f.sessionsRepo, f.hasher, s.logger, s.tracer)
			user, authToken, err := uc.SignUp(context.Background(), test.params)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if user != test.user {
				t.Errorf("\nExpected: %v\nGot: %v", test.user, user)
			}
			if authToken != test.authToken {
				t.Errorf("\nExpected: %v\nGot: %v", test.authToken, authToken)
			}
		})
	}
}

func (s *AuthUsecaseSuite) TestCheckAuth(t provider.T) {
	type fields struct {
		usersRepo    *usersMocks.MockRepository
		sessionsRepo *sessionsMocks.MockRepository
		userID       int
		authToken    string
	}

	type testCase struct {
		prepare   func(f *fields)
		userID    int
		authToken string
		err       error
	}

	tests := map[string]testCase{
		"normal": {
			prepare: func(f *fields) {
				gomock.InOrder(
					f.sessionsRepo.EXPECT().Get(gomock.Any(), f.userID, f.authToken).Return(f.userID, nil),
				)
			},
			userID:    21,
			authToken: "auth_token",
			err:       nil,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t provider.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				usersRepo:    usersMocks.NewMockRepository(ctrl),
				sessionsRepo: sessionsMocks.NewMockRepository(ctrl),
				userID:       test.userID,
				authToken:    test.authToken,
			}
			if test.prepare != nil {
				test.prepare(&f)
			}

			uc := authUsecase.New(f.usersRepo, f.sessionsRepo, hasherMocks.NewMockHasher(ctrl), s.logger, s.tracer)
			userID, err := uc.CheckAuth(context.Background(), test.userID, test.authToken)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if userID != test.userID {
				t.Errorf("\nExpected: %v\nGot: %v", test.userID, userID)
			}
		})
	}
}

func (s *AuthUsecaseSuite) TestLogout(t provider.T) {
	type fields struct {
		sessionsRepo *sessionsMocks.MockRepository
		userID       int
		authToken    string
	}

	type testCase struct {
		prepare   func(f *fields)
		userID    int
		authToken string
		err       error
	}

	tests := map[string]testCase{
		"normal": {
			prepare: func(f *fields) {
				f.sessionsRepo.EXPECT().Delete(gomock.Any(), f.userID, f.authToken).Return(nil)
			},
			userID:    21,
			authToken: "auth_token",
			err:       nil,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t provider.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				sessionsRepo: sessionsMocks.NewMockRepository(ctrl),
				userID:       test.userID,
				authToken:    test.authToken,
			}
			if test.prepare != nil {
				test.prepare(&f)
			}

			uc := authUsecase.New(usersMocks.NewMockRepository(ctrl), f.sessionsRepo, hasherMocks.NewMockHasher(ctrl),
				s.logger, s.tracer)
			err := uc.Logout(context.Background(), test.userID, test.authToken)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
		})
	}
}

func TestSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(AuthUsecaseSuite))
}
