package integration

import (
	"database/sql"
	"github.com/SlavaShagalov/my-trello-backend/internal/models"
	"github.com/SlavaShagalov/my-trello-backend/internal/pkg/config"
	pkgErrors "github.com/SlavaShagalov/my-trello-backend/internal/pkg/errors"
	pkgZap "github.com/SlavaShagalov/my-trello-backend/internal/pkg/log/zap"
	pkgDb "github.com/SlavaShagalov/my-trello-backend/internal/pkg/storages/postgres"
	pkgWorkspaces "github.com/SlavaShagalov/my-trello-backend/internal/workspaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"log"
	"os"
	"testing"

	workspacesRepo "github.com/SlavaShagalov/my-trello-backend/internal/workspaces/repository/postgres"
	workspacesUC "github.com/SlavaShagalov/my-trello-backend/internal/workspaces/usecase"
)

type WorkspacesSuite struct {
	suite.Suite
	db      *sql.DB
	logger  *zap.Logger
	logfile *os.File
	uc      pkgWorkspaces.Usecase
}

func (s *WorkspacesSuite) SetupSuite() {
	var err error
	s.logger, s.logfile, err = pkgZap.NewTestLogger("/logs/workspaces.log")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	config.SetTestPostgresConfig()
	s.db, err = pkgDb.NewStd(s.logger)
	s.Require().NoError(err)

	repo := workspacesRepo.New(s.db, s.logger)
	s.uc = workspacesUC.New(repo)
}

func (s *WorkspacesSuite) TearDownSuite() {
	err := s.db.Close()
	s.Require().NoError(err)

	err = s.logger.Sync()
	if err != nil {
		log.Println(err)
	}
	err = s.logfile.Close()
	if err != nil {
		log.Println(err)
	}
}

func (s *WorkspacesSuite) TestCreate() {
	type testCase struct {
		params *pkgWorkspaces.CreateParams
		err    error
	}

	tests := map[string]testCase{
		"normal": {
			params: &pkgWorkspaces.CreateParams{
				Title:       "University",
				Description: "BMSTU workspace",
				UserID:      3,
			},
			err: nil,
		},
		"user not found": {
			params: &pkgWorkspaces.CreateParams{
				Title:       "University",
				Description: "BMSTU workspace",
				UserID:      999,
			},
			err: pkgErrors.ErrUserNotFound,
		},
	}

	for name, test := range tests {
		s.Run(name, func() {
			workspace, err := s.uc.Create(test.params)

			assert.ErrorIs(s.T(), err, test.err, "unexpected error")

			if err == nil {
				assert.Equal(s.T(), test.params.UserID, workspace.UserID, "incorrect UserID")
				assert.Equal(s.T(), test.params.Title, workspace.Title, "incorrect Title")
				assert.Equal(s.T(), test.params.Description, workspace.Description, "incorrect Description")

				getWorkspace, err := s.uc.Get(workspace.ID)
				assert.NoError(s.T(), err, "failed to fetch workspace from the database")
				assert.Equal(s.T(), workspace.ID, getWorkspace.ID, "incorrect workspaceID")
				assert.Equal(s.T(), test.params.UserID, getWorkspace.UserID, "incorrect UserID")
				assert.Equal(s.T(), test.params.Title, getWorkspace.Title, "incorrect Title")
				assert.Equal(s.T(), test.params.Description, getWorkspace.Description, "incorrect Description")

				err = s.uc.Delete(workspace.ID)
				assert.NoError(s.T(), err, "failed to delete created workspace")
			}
		})
	}
}

func (s *WorkspacesSuite) TestList() {
	type testCase struct {
		userID     int
		workspaces []models.Workspace
		err        error
	}

	tests := map[string]testCase{
		"normal": {
			userID: 2,
			workspaces: []models.Workspace{
				{
					ID:          2,
					UserID:      2,
					Title:       "Проект \"Бета\"",
					Description: "Анализ данных и создание дашборда для маркетинговой отчетности",
				},
				{
					ID:          5,
					UserID:      2,
					Title:       "Проект \"Зета\"",
					Description: "Разработка API для интеграции с внешними сервисами",
				},
				{
					ID:          8,
					UserID:      2,
					Title:       "Проект \"Каппа\"",
					Description: "Создание руководств пользователя и документации",
				},
			},
			err: nil,
		},
		"empty result": {
			userID:     4,
			workspaces: []models.Workspace{},
			err:        nil,
		},
	}

	for name, test := range tests {
		s.Run(name, func() {
			workspaces, err := s.uc.List(test.userID)

			assert.ErrorIs(s.T(), err, test.err, "unexpected error")

			if err == nil {
				assert.Equal(s.T(), len(test.workspaces), len(workspaces), "incorrect workspaces length")
				for i := 0; i < len(test.workspaces); i++ {
					assert.Equal(s.T(), test.workspaces[i].ID, workspaces[i].ID, "incorrect workspaceID")
					assert.Equal(s.T(), test.workspaces[i].UserID, workspaces[i].UserID, "incorrect UserID")
					assert.Equal(s.T(), test.workspaces[i].Title, workspaces[i].Title, "incorrect Title")
					assert.Equal(s.T(), test.workspaces[i].Description, workspaces[i].Description, "incorrect Description")
				}
			}
		})
	}
}

func (s *WorkspacesSuite) TestGet() {
	type testCase struct {
		workspaceID int
		workspace   models.Workspace
		err         error
	}

	tests := map[string]testCase{
		"normal": {
			workspaceID: 8,
			workspace: models.Workspace{
				ID:          8,
				UserID:      2,
				Title:       "Проект \"Каппа\"",
				Description: "Создание руководств пользователя и документации",
			},
			err: nil,
		},
		"workspace not found": {
			workspaceID: 999,
			workspace:   models.Workspace{},
			err:         pkgErrors.ErrWorkspaceNotFound,
		},
	}

	for name, test := range tests {
		s.Run(name, func() {
			workspace, err := s.uc.Get(test.workspaceID)

			assert.ErrorIs(s.T(), err, test.err, "unexpected error")

			if err == nil {
				assert.Equal(s.T(), test.workspace.ID, workspace.ID, "incorrect workspaceID")
				assert.Equal(s.T(), test.workspace.UserID, workspace.UserID, "incorrect UserID")
				assert.Equal(s.T(), test.workspace.Title, workspace.Title, "incorrect Title")
				assert.Equal(s.T(), test.workspace.Description, workspace.Description, "incorrect Description")
			}
		})
	}
}

func (s *WorkspacesSuite) TestFullUpdate() {
	type testCase struct {
		params *pkgWorkspaces.FullUpdateParams
		err    error
	}

	tests := map[string]testCase{
		"normal": {
			params: &pkgWorkspaces.FullUpdateParams{
				Title:       "University",
				Description: "BMSTU workspace",
			},
			err: nil,
		},
	}

	for name, test := range tests {
		s.Run(name, func() {
			tempWorkspace, err := s.uc.Create(&pkgWorkspaces.CreateParams{
				Title:       "Temp Workspace",
				Description: "Temp Workspace Description",
				UserID:      2,
			})
			require.NoError(s.T(), err, "failed to create temp workspace")

			test.params.ID = tempWorkspace.ID
			workspace, err := s.uc.FullUpdate(test.params)
			assert.ErrorIs(s.T(), err, test.err, "unexpected error")

			if err == nil {
				// check updated workspace
				assert.Equal(s.T(), test.params.ID, workspace.ID, "incorrect workspaceID")
				assert.Equal(s.T(), test.params.Title, workspace.Title, "incorrect Title")
				assert.Equal(s.T(), test.params.Description, workspace.Description, "incorrect Description")

				// check workspace in storages
				getWorkspace, err := s.uc.Get(workspace.ID)
				assert.NoError(s.T(), err, "failed to fetch workspace from the database")
				assert.Equal(s.T(), workspace.ID, getWorkspace.ID, "incorrect workspaceID")
				assert.Equal(s.T(), workspace.UserID, getWorkspace.UserID, "incorrect UserID")
				assert.Equal(s.T(), test.params.Title, getWorkspace.Title, "incorrect Title")
				assert.Equal(s.T(), test.params.Description, getWorkspace.Description, "incorrect Description")
			}

			err = s.uc.Delete(tempWorkspace.ID)
			require.NoError(s.T(), err, "failed to delete temp workspace")
		})
	}
}

func (s *WorkspacesSuite) TestPartialUpdate() {
	type testCase struct {
		params    *pkgWorkspaces.PartialUpdateParams
		workspace models.Workspace
		err       error
	}

	tests := map[string]testCase{
		"full update": {
			params: &pkgWorkspaces.PartialUpdateParams{
				Title:             "University",
				UpdateTitle:       true,
				Description:       "BMSTU workspace",
				UpdateDescription: true,
			},
			workspace: models.Workspace{
				Title:       "University",
				Description: "BMSTU workspace",
			},
			err: nil,
		},
		"only title update": {
			params: &pkgWorkspaces.PartialUpdateParams{
				Title:       "New University",
				UpdateTitle: true,
			},
			workspace: models.Workspace{
				Title:       "New University",
				Description: "Temp Workspace Description",
			},
			err: nil,
		},
		"only description update": {
			params: &pkgWorkspaces.PartialUpdateParams{
				Description:       "New BMSTU workspace",
				UpdateDescription: true,
			},
			workspace: models.Workspace{
				Title:       "Temp Workspace",
				Description: "New BMSTU workspace",
			},
			err: nil,
		},
	}

	for name, test := range tests {
		s.Run(name, func() {
			tempWorkspace, err := s.uc.Create(&pkgWorkspaces.CreateParams{
				Title:       "Temp Workspace",
				Description: "Temp Workspace Description",
				UserID:      2,
			})
			require.NoError(s.T(), err, "failed to create temp workspace")

			test.params.ID = tempWorkspace.ID
			workspace, err := s.uc.PartialUpdate(test.params)
			assert.ErrorIs(s.T(), err, test.err, "unexpected error")

			if err == nil {
				// check updated workspace
				assert.Equal(s.T(), test.params.ID, workspace.ID, "incorrect workspaceID")
				assert.Equal(s.T(), test.workspace.Title, workspace.Title, "incorrect Title")
				assert.Equal(s.T(), test.workspace.Description, workspace.Description, "incorrect Description")

				// check workspace in storages
				getWorkspace, err := s.uc.Get(workspace.ID)
				assert.NoError(s.T(), err, "failed to fetch workspace from the database")
				assert.Equal(s.T(), workspace.UserID, getWorkspace.UserID, "incorrect UserID")
				assert.Equal(s.T(), test.workspace.Title, getWorkspace.Title, "incorrect Title")
				assert.Equal(s.T(), test.workspace.Description, getWorkspace.Description, "incorrect Description")
			}

			err = s.uc.Delete(tempWorkspace.ID)
			require.NoError(s.T(), err, "failed to delete temp workspace")
		})
	}
}

func (s *WorkspacesSuite) TestDelete() {
	type testCase struct {
		setupWorkspace func() (models.Workspace, error)
		err            error
	}

	tests := map[string]testCase{
		"normal": {
			setupWorkspace: func() (models.Workspace, error) {
				return s.uc.Create(&pkgWorkspaces.CreateParams{
					Title:       "Test Workspace",
					Description: "Test Workspace Description",
					UserID:      1,
				})
			},
			err: nil,
		},
		"workspace not found": {
			setupWorkspace: func() (models.Workspace, error) {
				return models.Workspace{ID: 999}, nil
			},
			err: pkgErrors.ErrWorkspaceNotFound,
		},
	}

	for name, test := range tests {
		s.Run(name, func() {
			workspace, err := test.setupWorkspace()
			s.Require().NoError(err)

			err = s.uc.Delete(workspace.ID)
			assert.ErrorIs(s.T(), err, test.err, "unexpected error")

			if test.err == nil {
				_, err = s.uc.Get(workspace.ID)
				assert.ErrorIs(s.T(), err, pkgErrors.ErrWorkspaceNotFound, "workspace should be deleted")
			}
		})
	}
}

func TestWorkspaceSuite(t *testing.T) {
	suite.Run(t, new(WorkspacesSuite))
}
