package workspaces

import (
	pkgZap "github.com/SlavaShagalov/my-trello-backend/internal/pkg/log/zap"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"

	"github.com/SlavaShagalov/my-trello-backend/internal/models"
	pkgErrors "github.com/SlavaShagalov/my-trello-backend/internal/pkg/errors"
	pkgWorkspaces "github.com/SlavaShagalov/my-trello-backend/internal/workspaces"
	"github.com/SlavaShagalov/my-trello-backend/internal/workspaces/mocks"
	workspacesUsecase "github.com/SlavaShagalov/my-trello-backend/internal/workspaces/usecase"
)

type WorkspacesUsecaseSuite struct {
	suite.Suite
	logger *zap.Logger
}

func (s *WorkspacesUsecaseSuite) BeforeAll(t provider.T) {
	t.WithNewStep("SetupSuite step", func(ctx provider.StepCtx) {})

	s.logger = pkgZap.NewTestLogger()
}

func (s *WorkspacesUsecaseSuite) AfterAll(t provider.T) {
	t.WithNewStep("TearDownSuite step", func(ctx provider.StepCtx) {})

	_ = s.logger.Sync()
}

func (s *WorkspacesUsecaseSuite) BeforeEach(t provider.T) {
	t.WithNewStep("SetupTest step", func(ctx provider.StepCtx) {})
}

func (s *WorkspacesUsecaseSuite) AfterEach(t provider.T) {
	t.WithNewStep("TearDownTest step", func(ctx provider.StepCtx) {})
}

func (s *WorkspacesUsecaseSuite) TestCreate(t provider.T) {
	type fields struct {
		repo      *mocks.MockRepository
		params    *pkgWorkspaces.CreateParams
		workspace *models.Workspace
	}

	type testCase struct {
		prepare   func(f *fields)
		params    *pkgWorkspaces.CreateParams
		workspace models.Workspace
		err       error
	}

	tests := map[string]testCase{
		"normal": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Create(f.params).Return(*f.workspace, nil)
			},
			params:    &pkgWorkspaces.CreateParams{Title: "University", Description: "BMSTU workspace", UserID: 27},
			workspace: models.Workspace{ID: 21, UserID: 27, Title: "University", Description: "BMSTU workspace"},
			err:       nil,
		},
		"user not found": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Create(f.params).Return(*f.workspace, pkgErrors.ErrUserNotFound)
			},
			params:    &pkgWorkspaces.CreateParams{Title: "University", Description: "BMSTU workspace", UserID: 27},
			workspace: models.Workspace{},
			err:       pkgErrors.ErrUserNotFound,
		},
		"db error": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Create(f.params).Return(*f.workspace, pkgErrors.ErrDb)
			},
			params:    &pkgWorkspaces.CreateParams{Title: "University", Description: "BMSTU workspace", UserID: 27},
			workspace: models.Workspace{},
			err:       pkgErrors.ErrDb,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t provider.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl), params: test.params, workspace: &test.workspace}
			if test.prepare != nil {
				test.prepare(&f)
			}

			uc := workspacesUsecase.New(f.repo)
			workspace, err := uc.Create(test.params)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if workspace != test.workspace {
				t.Errorf("\nExpected: %v\nGot: %v", test.workspace, workspace)
			}
		})
	}
}

func (s *WorkspacesUsecaseSuite) TestList(t provider.T) {
	type fields struct {
		repo       *mocks.MockRepository
		userID     int
		workspaces []models.Workspace
	}

	type testCase struct {
		prepare    func(f *fields)
		userID     int
		workspaces []models.Workspace
		err        error
	}

	tests := map[string]testCase{
		"normal": {
			prepare: func(f *fields) {
				f.repo.EXPECT().List(f.userID).Return(f.workspaces, nil)
			},
			userID: 27,
			workspaces: []models.Workspace{
				{ID: 21, UserID: 27, Title: "University", Description: "BMSTU workspace"},
				{ID: 22, UserID: 27, Title: "Work", Description: "Work workspace"},
				{ID: 23, UserID: 27, Title: "Life", Description: "Life workspace"},
			},
			err: nil,
		},
		"empty result": {
			prepare: func(f *fields) {
				f.repo.EXPECT().List(f.userID).Return(f.workspaces, nil)
			},
			userID:     27,
			workspaces: []models.Workspace{},
			err:        nil,
		},
		"user not found": {
			prepare: func(f *fields) {
				f.repo.EXPECT().List(f.userID).Return(f.workspaces, pkgErrors.ErrUserNotFound)
			},
			userID:     27,
			workspaces: nil,
			err:        pkgErrors.ErrUserNotFound,
		},
		"db error": {
			prepare: func(f *fields) {
				f.repo.EXPECT().List(f.userID).Return(f.workspaces, pkgErrors.ErrDb)
			},
			userID:     27,
			workspaces: nil,
			err:        pkgErrors.ErrDb,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t provider.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl), userID: test.userID, workspaces: test.workspaces}
			if test.prepare != nil {
				test.prepare(&f)
			}

			uc := workspacesUsecase.New(f.repo)
			workspaces, err := uc.List(test.userID)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if !reflect.DeepEqual(workspaces, test.workspaces) {
				t.Errorf("\nExpected: %v\nGot: %v", test.workspaces, workspaces)
			}
		})
	}
}

func (s *WorkspacesUsecaseSuite) TestGet(t provider.T) {
	type fields struct {
		repo        *mocks.MockRepository
		workspaceID int
		workspace   *models.Workspace
	}

	type testCase struct {
		prepare     func(f *fields)
		workspaceID int
		workspace   models.Workspace
		err         error
	}

	tests := map[string]testCase{
		"normal": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Get(f.workspaceID).Return(*f.workspace, nil)
			},
			workspaceID: 21,
			workspace:   models.Workspace{ID: 21, UserID: 27, Title: "University", Description: "BMSTU workspace"},
			err:         nil,
		},
		"user not found": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Get(f.workspaceID).Return(*f.workspace, pkgErrors.ErrUserNotFound)
			},
			workspaceID: 21,
			workspace:   models.Workspace{},
			err:         pkgErrors.ErrUserNotFound,
		},
		"db error": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Get(f.workspaceID).Return(*f.workspace, pkgErrors.ErrDb)
			},
			workspaceID: 21,
			workspace:   models.Workspace{},
			err:         pkgErrors.ErrDb,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t provider.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl), workspaceID: test.workspaceID, workspace: &test.workspace}
			if test.prepare != nil {
				test.prepare(&f)
			}

			uc := workspacesUsecase.New(f.repo)
			workspace, err := uc.Get(test.workspaceID)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if workspace != test.workspace {
				t.Errorf("\nExpected: %v\nGot: %v", test.workspace, workspace)
			}
		})
	}
}

func (s *WorkspacesUsecaseSuite) TestFullUpdate(t provider.T) {
	type fields struct {
		repo      *mocks.MockRepository
		params    *pkgWorkspaces.FullUpdateParams
		workspace *models.Workspace
	}

	type testCase struct {
		prepare   func(f *fields)
		params    *pkgWorkspaces.FullUpdateParams
		workspace models.Workspace
		err       error
	}

	tests := map[string]testCase{
		"normal": {
			prepare: func(f *fields) {
				f.repo.EXPECT().FullUpdate(f.params).Return(*f.workspace, nil)
			},
			params:    &pkgWorkspaces.FullUpdateParams{ID: 21, Title: "University", Description: "BMSTU workspace"},
			workspace: models.Workspace{ID: 21, UserID: 27, Title: "University", Description: "BMSTU workspace"},
			err:       nil,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t provider.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl), params: test.params, workspace: &test.workspace}
			if test.prepare != nil {
				test.prepare(&f)
			}

			uc := workspacesUsecase.New(f.repo)
			workspace, err := uc.FullUpdate(test.params)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if workspace != test.workspace {
				t.Errorf("\nExpected: %v\nGot: %v", test.workspace, workspace)
			}
		})
	}
}

func (s *WorkspacesUsecaseSuite) TestPartialUpdate(t provider.T) {
	type fields struct {
		repo      *mocks.MockRepository
		params    *pkgWorkspaces.PartialUpdateParams
		workspace *models.Workspace
	}

	type testCase struct {
		prepare   func(f *fields)
		params    *pkgWorkspaces.PartialUpdateParams
		workspace models.Workspace
		err       error
	}

	tests := map[string]testCase{
		"normal": {
			prepare: func(f *fields) {
				f.repo.EXPECT().PartialUpdate(f.params).Return(*f.workspace, nil)
			},
			params: &pkgWorkspaces.PartialUpdateParams{
				ID:                21,
				Title:             "University",
				UpdateTitle:       true,
				Description:       "BMSTU workspace",
				UpdateDescription: true,
			},
			workspace: models.Workspace{ID: 21, UserID: 27, Title: "University", Description: "BMSTU workspace"},
			err:       nil,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t provider.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl), params: test.params, workspace: &test.workspace}
			if test.prepare != nil {
				test.prepare(&f)
			}

			uc := workspacesUsecase.New(f.repo)
			workspace, err := uc.PartialUpdate(test.params)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if workspace != test.workspace {
				t.Errorf("\nExpected: %v\nGot: %v", test.workspace, workspace)
			}
		})
	}
}

func (s *WorkspacesUsecaseSuite) TestDelete(t provider.T) {
	type fields struct {
		repo        *mocks.MockRepository
		workspaceID int
	}

	type testCase struct {
		prepare     func(f *fields)
		workspaceID int
		err         error
	}

	tests := map[string]testCase{
		"normal": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Delete(f.workspaceID).Return(nil)
			},
			workspaceID: 21,
			err:         nil,
		},
		"workspace not found": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Delete(f.workspaceID).Return(pkgErrors.ErrWorkspaceNotFound)
			},
			workspaceID: 21,
			err:         pkgErrors.ErrWorkspaceNotFound,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t provider.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl), workspaceID: test.workspaceID}
			if test.prepare != nil {
				test.prepare(&f)
			}

			uc := workspacesUsecase.New(f.repo)
			err := uc.Delete(test.workspaceID)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
		})
	}
}

func TestSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(WorkspacesUsecaseSuite))
}
