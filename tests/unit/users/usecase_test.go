package users

import (
	"github.com/SlavaShagalov/my-trello-backend/internal/pkg/config"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/pkg/errors"

	imgMocks "github.com/SlavaShagalov/my-trello-backend/internal/images/mocks"
	"github.com/SlavaShagalov/my-trello-backend/internal/models"
	pkgErrors "github.com/SlavaShagalov/my-trello-backend/internal/pkg/errors"
	pkgUsers "github.com/SlavaShagalov/my-trello-backend/internal/users"
	"github.com/SlavaShagalov/my-trello-backend/internal/users/mocks"
	usersUsecase "github.com/SlavaShagalov/my-trello-backend/internal/users/usecase"
)

type UsersUsecaseSuite struct {
	suite.Suite
}

func (s *UsersUsecaseSuite) BeforeAll(t provider.T) {
	t.WithNewStep("SetupSuite step", func(ctx provider.StepCtx) {})

	config.SetDefaultValidationConfig()
}

func (s *UsersUsecaseSuite) AfterAll(t provider.T) {
	t.WithNewStep("TearDownSuite step", func(ctx provider.StepCtx) {})
}

func (s *UsersUsecaseSuite) BeforeEach(t provider.T) {
	t.WithNewStep("SetupTest step", func(ctx provider.StepCtx) {})
}

func (s *UsersUsecaseSuite) AfterEach(t provider.T) {
	t.WithNewStep("TearDownTest step", func(ctx provider.StepCtx) {})
}

func (s *UsersUsecaseSuite) TestList(t provider.T) {
	type fields struct {
		repo    *mocks.MockRepository
		imgRepo *imgMocks.MockRepository
		users   []models.User
	}

	type testCase struct {
		prepare func(f *fields)
		users   []models.User
		err     error
	}

	tests := map[string]testCase{
		"normal": {
			prepare: func(f *fields) {
				f.repo.EXPECT().List().Return(f.users, nil)
			},
			users: []models.User{
				{
					ID:       21,
					Username: "slava",
					Password: "hash1",
					Email:    "slava@vk.com",
					Name:     "Slava",
				},
				{
					ID:       22,
					Username: "petya",
					Password: "hash2",
					Email:    "petya@vk.com",
					Name:     "Petya",
				},
				{
					ID:       23,
					Username: "misha",
					Password: "hash3",
					Email:    "misha@vk.com",
					Name:     "Misha",
				},
			},
			err: nil,
		},
		"empty result": {
			prepare: func(f *fields) {
				f.repo.EXPECT().List().Return(f.users, nil)
			},
			users: []models.User{},
			err:   nil,
		},
		"storages error": {
			prepare: func(f *fields) {
				f.repo.EXPECT().List().Return(f.users, pkgErrors.ErrDb)
			},
			users: nil,
			err:   pkgErrors.ErrDb,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t provider.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				repo:    mocks.NewMockRepository(ctrl),
				imgRepo: imgMocks.NewMockRepository(ctrl),
				users:   test.users,
			}
			if test.prepare != nil {
				test.prepare(&f)
			}

			uc := usersUsecase.New(f.repo, f.imgRepo)
			workspaces, err := uc.List()
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if !reflect.DeepEqual(workspaces, test.users) {
				t.Errorf("\nExpected: %v\nGot: %v", test.users, workspaces)
			}
		})
	}
}

func (s *UsersUsecaseSuite) TestGet(t provider.T) {
	type fields struct {
		repo    *mocks.MockRepository
		imgRepo *imgMocks.MockRepository
		userID  int
		user    *models.User
	}

	type testCase struct {
		prepare func(f *fields)
		userID  int
		user    models.User
		err     error
	}

	tests := map[string]testCase{
		"normal": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Get(f.userID).Return(*f.user, nil)
			},
			userID: 21,
			user: models.User{
				ID:       21,
				Username: "slava",
				Password: "hash1",
				Email:    "slava@vk.com",
				Name:     "Slava",
			},
			err: nil,
		},
		"storages error": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Get(f.userID).Return(*f.user, pkgErrors.ErrDb)
			},
			userID: 21,
			user:   models.User{},
			err:    pkgErrors.ErrDb,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t provider.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl), userID: test.userID, user: &test.user}
			if test.prepare != nil {
				test.prepare(&f)
			}

			uc := usersUsecase.New(f.repo, f.imgRepo)
			user, err := uc.Get(test.userID)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if user != test.user {
				t.Errorf("\nExpected: %v\nGot: %v", test.user, user)
			}
		})
	}
}

func (s *UsersUsecaseSuite) TestGetByUsername(t provider.T) {
	type fields struct {
		repo     *mocks.MockRepository
		imgRepo  *imgMocks.MockRepository
		username string
		user     *models.User
	}

	type testCase struct {
		prepare  func(f *fields)
		username string
		user     models.User
		err      error
	}

	tests := map[string]testCase{
		"normal": {
			prepare: func(f *fields) {
				f.repo.EXPECT().GetByUsername(f.username).Return(*f.user, nil)
			},
			username: "slava",
			user: models.User{
				ID:       21,
				Username: "slava",
				Password: "hash1",
				Email:    "slava@vk.com",
				Name:     "Slava",
			},
			err: nil,
		},
		"storages error": {
			prepare: func(f *fields) {
				f.repo.EXPECT().GetByUsername(f.username).Return(*f.user, pkgErrors.ErrDb)
			},
			username: "slava",
			user:     models.User{},
			err:      pkgErrors.ErrDb,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t provider.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl), username: test.username, user: &test.user}
			if test.prepare != nil {
				test.prepare(&f)
			}

			uc := usersUsecase.New(f.repo, f.imgRepo)
			user, err := uc.GetByUsername(test.username)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if user != test.user {
				t.Errorf("\nExpected: %v\nGot: %v", test.user, user)
			}
		})
	}
}

func (s *UsersUsecaseSuite) TestFullUpdate(t provider.T) {
	type fields struct {
		repo    *mocks.MockRepository
		imgRepo *imgMocks.MockRepository
		params  *pkgUsers.FullUpdateParams
		user    *models.User
	}

	type testCase struct {
		prepare func(f *fields)
		params  *pkgUsers.FullUpdateParams
		user    models.User
		err     error
	}

	tests := map[string]testCase{
		"normal": {
			prepare: func(f *fields) {
				f.repo.EXPECT().FullUpdate(f.params).Return(*f.user, nil)
				f.repo.EXPECT().GetByUsername(f.params.Username).Return(*f.user, nil)
			},
			params: &pkgUsers.FullUpdateParams{
				ID:       21,
				Username: "slava",
				Email:    "slava@vk.com",
				Name:     "Slava",
			},
			user: models.User{ID: 21, Username: "slava", Email: "slava@vk.com", Name: "Slava"},
			err:  nil,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t provider.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl), params: test.params, user: &test.user}
			if test.prepare != nil {
				test.prepare(&f)
			}

			uc := usersUsecase.New(f.repo, f.imgRepo)
			user, err := uc.FullUpdate(test.params)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if user != test.user {
				t.Errorf("\nExpected: %v\nGot: %v", test.user, user)
			}
		})
	}
}

func (s *UsersUsecaseSuite) TestPartialUpdate(t provider.T) {
	type fields struct {
		repo    *mocks.MockRepository
		imgRepo *imgMocks.MockRepository
		params  *pkgUsers.PartialUpdateParams
		user    *models.User
	}

	type testCase struct {
		prepare func(f *fields)
		params  *pkgUsers.PartialUpdateParams
		user    models.User
		err     error
	}

	tests := map[string]testCase{
		"normal": {
			prepare: func(f *fields) {
				f.repo.EXPECT().PartialUpdate(f.params).Return(*f.user, nil)
				f.repo.EXPECT().GetByUsername(f.params.Username).Return(*f.user, nil)
			},
			params: &pkgUsers.PartialUpdateParams{
				ID:             21,
				Username:       "slava",
				UpdateUsername: true,
				Email:          "slava@vk.com",
				UpdateEmail:    true,
				Name:           "Slava",
				UpdateName:     true,
			},
			user: models.User{ID: 21, Username: "slava", Email: "slava@vk.com", Name: "Slava"},
			err:  nil,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t provider.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl), params: test.params, user: &test.user}
			if test.prepare != nil {
				test.prepare(&f)
			}

			uc := usersUsecase.New(f.repo, f.imgRepo)
			user, err := uc.PartialUpdate(test.params)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if user != test.user {
				t.Errorf("\nExpected: %v\nGot: %v", test.user, user)
			}
		})
	}
}

func (s *UsersUsecaseSuite) TestDelete(t provider.T) {
	type fields struct {
		repo    *mocks.MockRepository
		imgRepo *imgMocks.MockRepository
		userID  int
	}

	type testCase struct {
		prepare func(f *fields)
		userID  int
		err     error
	}

	tests := map[string]testCase{
		"normal": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Delete(f.userID).Return(nil)
			},
			userID: 21,
			err:    nil,
		},
		"user not found": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Delete(f.userID).Return(pkgErrors.ErrUserNotFound)
			},
			userID: 21,
			err:    pkgErrors.ErrUserNotFound,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t provider.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl), userID: test.userID}
			if test.prepare != nil {
				test.prepare(&f)
			}

			uc := usersUsecase.New(f.repo, f.imgRepo)
			err := uc.Delete(test.userID)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
		})
	}
}

func TestSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(UsersUsecaseSuite))
}
