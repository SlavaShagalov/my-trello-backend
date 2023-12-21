package workspaces

import (
	pkgZap "github.com/SlavaShagalov/my-trello-backend/internal/pkg/log/zap"
	"github.com/SlavaShagalov/my-trello-backend/tests/utils/builder"
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
	logger    *zap.Logger
	wsBuilder *builder.WorkspaceBuilder
}

func (s *WorkspacesUsecaseSuite) BeforeAll(t provider.T) {
	t.WithNewStep("SetupSuite step", func(ctx provider.StepCtx) {})

	s.logger = pkgZap.NewDevelopLogger()
}

func (s *WorkspacesUsecaseSuite) AfterAll(t provider.T) {
	t.WithNewStep("TearDownSuite step", func(ctx provider.StepCtx) {})

	_ = s.logger.Sync()
}

func (s *WorkspacesUsecaseSuite) BeforeEach(t provider.T) {
	t.WithNewStep("SetupTest step", func(ctx provider.StepCtx) {})

	s.wsBuilder = builder.NewWorkspaceBuilder()
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
			params: &pkgWorkspaces.CreateParams{Title: "University", Description: "BMSTU workspace", UserID: 27},
			workspace: s.wsBuilder.
				WithID(21).
				WithUserID(27).
				WithTitle("University").
				WithDescription("BMSTU workspace").
				Build(),
			err: nil,
		},
		"user not found": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Create(f.params).Return(*f.workspace, pkgErrors.ErrUserNotFound)
			},
			params:    &pkgWorkspaces.CreateParams{Title: "University", Description: "BMSTU workspace", UserID: 27},
			workspace: s.wsBuilder.Build(),
			err:       pkgErrors.ErrUserNotFound,
		},
		"db error": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Create(f.params).Return(*f.workspace, pkgErrors.ErrDb)
			},
			params:    &pkgWorkspaces.CreateParams{Title: "University", Description: "BMSTU workspace", UserID: 27},
			workspace: s.wsBuilder.Build(),
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
				s.wsBuilder.
					WithID(21).
					WithUserID(27).
					WithTitle("University").
					WithDescription("BMSTU workspace").
					Build(),
				s.wsBuilder.
					WithID(22).
					WithUserID(27).
					WithTitle("Work").
					WithDescription("Work workspace").
					Build(),
				s.wsBuilder.
					WithID(23).
					WithUserID(27).
					WithTitle("Life").
					WithDescription("Life workspace").
					Build(),
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
			workspace: s.wsBuilder.
				WithID(21).
				WithUserID(27).
				WithTitle("University").
				WithDescription("BMSTU workspace").
				Build(),
			err: nil,
		},
		"user not found": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Get(f.workspaceID).Return(*f.workspace, pkgErrors.ErrUserNotFound)
			},
			workspaceID: 21,
			workspace:   s.wsBuilder.Build(),
			err:         pkgErrors.ErrUserNotFound,
		},
		"db error": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Get(f.workspaceID).Return(*f.workspace, pkgErrors.ErrDb)
			},
			workspaceID: 21,
			workspace:   s.wsBuilder.Build(),
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
			params: &pkgWorkspaces.FullUpdateParams{ID: 21, Title: "University", Description: "BMSTU workspace"},
			workspace: s.wsBuilder.
				WithID(21).
				WithUserID(27).
				WithTitle("University").
				WithDescription("BMSTU workspace").
				Build(),
			err: nil,
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
			workspace: s.wsBuilder.
				WithID(21).
				WithUserID(27).
				WithTitle("University").
				WithDescription("BMSTU workspace").
				Build(),
			err: nil,
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
