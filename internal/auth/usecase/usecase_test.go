package usecase

import (
	pkgAuth "github.com/SlavaShagalov/my-trello-backend/internal/auth"
	"github.com/SlavaShagalov/my-trello-backend/internal/models"
	pkgErrors "github.com/SlavaShagalov/my-trello-backend/internal/pkg/errors"
	"github.com/SlavaShagalov/my-trello-backend/internal/users"
	"github.com/pkg/errors"

	hasherMocks "github.com/SlavaShagalov/my-trello-backend/internal/pkg/hasher/mocks"
	sessionsMocks "github.com/SlavaShagalov/my-trello-backend/internal/sessions/mocks"
	usersMocks "github.com/SlavaShagalov/my-trello-backend/internal/users/mocks"

	"github.com/golang/mock/gomock"
	"testing"
)

func TestUsecase_SignIn(t *testing.T) {
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
					f.usersRepo.EXPECT().GetByUsername(f.params.Username).Return(*f.user, nil),
					f.hasher.EXPECT().CompareHashAndPassword(f.user.Password, f.params.Password).Return(nil),
					f.sessionsRepo.EXPECT().Create(f.user.ID).Return(f.authToken, nil),
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
		t.Run(name, func(t *testing.T) {
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

			uc := New(f.usersRepo, f.sessionsRepo, f.hasher)
			user, authToken, err := uc.SignIn(test.params)
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

func TestUsecase_SignUp(t *testing.T) {
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
					f.usersRepo.EXPECT().GetByUsername(f.params.Username).
						Return(models.User{}, pkgErrors.ErrUserNotFound),
					f.hasher.EXPECT().GetHashedPassword(f.params.Password).Return(f.user.Password, nil),
					f.usersRepo.EXPECT().Create(&users.CreateParams{
						Name:           f.params.Name,
						Username:       f.params.Username,
						Email:          f.params.Email,
						HashedPassword: f.user.Password,
					}).Return(*f.user, nil),
					f.sessionsRepo.EXPECT().Create(f.user.ID).Return(f.authToken, nil),
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
		t.Run(name, func(t *testing.T) {
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

			uc := New(f.usersRepo, f.sessionsRepo, f.hasher)
			user, authToken, err := uc.SignUp(test.params)
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

func TestUsecase_CheckAuth(t *testing.T) {
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
					f.sessionsRepo.EXPECT().Get(f.userID, f.authToken).Return(f.userID, nil),
				)
			},
			userID:    21,
			authToken: "auth_token",
			err:       nil,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
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

			uc := New(f.usersRepo, f.sessionsRepo, hasherMocks.NewMockHasher(ctrl))
			userID, err := uc.CheckAuth(test.userID, test.authToken)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if userID != test.userID {
				t.Errorf("\nExpected: %v\nGot: %v", test.userID, userID)
			}
		})
	}
}

func TestUsecase_Logout(t *testing.T) {
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
				f.sessionsRepo.EXPECT().Delete(f.userID, f.authToken).Return(nil)
			},
			userID:    21,
			authToken: "auth_token",
			err:       nil,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
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

			uc := New(usersMocks.NewMockRepository(ctrl), f.sessionsRepo, hasherMocks.NewMockHasher(ctrl))
			err := uc.Logout(test.userID, test.authToken)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
		})
	}
}
