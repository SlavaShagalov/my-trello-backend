package integration

import (
	"database/sql"
	"github.com/SlavaShagalov/my-trello-backend/internal/models"
	"github.com/SlavaShagalov/my-trello-backend/internal/pkg/config"
	pkgErrors "github.com/SlavaShagalov/my-trello-backend/internal/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"testing"

	pkgLists "github.com/SlavaShagalov/my-trello-backend/internal/lists"
	pkgZap "github.com/SlavaShagalov/my-trello-backend/internal/pkg/log/zap"
	pkgDb "github.com/SlavaShagalov/my-trello-backend/internal/pkg/storages"

	listsRepo "github.com/SlavaShagalov/my-trello-backend/internal/lists/repository/postgres"
	listsUC "github.com/SlavaShagalov/my-trello-backend/internal/lists/usecase"
)

type ListsSuite struct {
	suite.Suite
	db     *sql.DB
	logger *zap.Logger
	uc     pkgLists.Usecase
}

func (s *ListsSuite) SetupSuite() {
	s.logger = pkgZap.NewTestLogger()

	var err error
	config.SetTestPostgresConfig()
	s.db, err = pkgDb.NewPostgres(s.logger)
	s.Require().NoError(err)

	repo := listsRepo.NewRepository(s.db, s.logger)
	s.uc = listsUC.NewUsecase(repo)
}

func (s *ListsSuite) TearDownSuite() {
	err := s.db.Close()
	s.Require().NoError(err)

	_ = s.logger.Sync()
}

func (s *ListsSuite) TestCreate() {
	type testCase struct {
		params *pkgLists.CreateParams
		err    error
	}

	tests := map[string]testCase{
		"normal": {
			params: &pkgLists.CreateParams{
				Title:   "MathStat",
				BoardID: 3,
			},
			err: nil,
		},
		"board not found": {
			params: &pkgLists.CreateParams{
				Title:   "MathStat",
				BoardID: 999,
			},
			err: pkgErrors.ErrBoardNotFound,
		},
	}

	for name, test := range tests {
		s.Run(name, func() {
			list, err := s.uc.Create(test.params)
			assert.ErrorIs(s.T(), err, test.err, "unexpected error")

			if err == nil {
				assert.Equal(s.T(), test.params.BoardID, list.BoardID, "incorrect BoardID")
				assert.Equal(s.T(), test.params.Title, list.Title, "incorrect Title")

				getList, err := s.uc.Get(list.ID)
				assert.NoError(s.T(), err, "failed to fetch list from the database")
				assert.Equal(s.T(), list.ID, getList.ID, "incorrect listID")
				assert.Equal(s.T(), test.params.BoardID, getList.BoardID, "incorrect BoardID")
				assert.Equal(s.T(), test.params.Title, getList.Title, "incorrect Title")

				err = s.uc.Delete(list.ID)
				assert.NoError(s.T(), err, "failed to delete created list")
			}
		})
	}
}

func (s *ListsSuite) TestList() {
	type testCase struct {
		boardID int
		lists   []models.List
		err     error
	}

	tests := map[string]testCase{
		"normal": {
			boardID: 2,
			lists: []models.List{
				{
					ID:       4,
					BoardID:  2,
					Title:    "Анализ данных",
					Position: 1,
				},
				{
					ID:       5,
					BoardID:  2,
					Title:    "Подготовка отчета",
					Position: 2,
				},
				{
					ID:       6,
					BoardID:  2,
					Title:    "Маркетинговые мероприятия",
					Position: 3,
				},
			},
			err: nil,
		},
		"empty result": {
			boardID: 11,
			lists:   []models.List{},
			err:     nil,
		},
	}

	for name, test := range tests {
		s.Run(name, func() {
			lists, err := s.uc.ListByBoard(test.boardID)

			assert.ErrorIs(s.T(), err, test.err, "unexpected error")

			if err == nil {
				assert.Equal(s.T(), len(test.lists), len(lists), "incorrect lists length")
				for i := 0; i < len(test.lists); i++ {
					assert.Equal(s.T(), test.lists[i].ID, lists[i].ID, "incorrect listID")
					assert.Equal(s.T(), test.lists[i].BoardID, lists[i].BoardID, "incorrect BoardID")
					assert.Equal(s.T(), test.lists[i].Title, lists[i].Title, "incorrect Title")
					assert.Equal(s.T(), test.lists[i].Position, lists[i].Position, "incorrect Position")
				}
			}
		})
	}
}

func (s *ListsSuite) TestGet() {
	type testCase struct {
		listID int
		list   models.List
		err    error
	}

	tests := map[string]testCase{
		"normal": {
			listID: 8,
			list: models.List{
				ID:       8,
				BoardID:  3,
				Title:    "Прототипирование",
				Position: 2,
			},
			err: nil,
		},
		"list not found": {
			listID: 999,
			list:   models.List{},
			err:    pkgErrors.ErrListNotFound,
		},
	}

	for name, test := range tests {
		s.Run(name, func() {
			list, err := s.uc.Get(test.listID)

			assert.ErrorIs(s.T(), err, test.err, "unexpected error")

			if err == nil {
				assert.Equal(s.T(), test.list.ID, list.ID, "incorrect listID")
				assert.Equal(s.T(), test.list.BoardID, list.BoardID, "incorrect BoardID")
				assert.Equal(s.T(), test.list.Title, list.Title, "incorrect Title")
				assert.Equal(s.T(), test.list.Position, list.Position, "incorrect Position")
			}
		})
	}
}

func (s *ListsSuite) TestFullUpdate() {
	type testCase struct {
		params *pkgLists.FullUpdateParams
		err    error
	}

	tests := map[string]testCase{
		"normal": {
			params: &pkgLists.FullUpdateParams{
				Title:   "MathStat",
				BoardID: 3,
			},
			err: nil,
		},
	}

	for name, test := range tests {
		s.Run(name, func() {
			tempList, err := s.uc.Create(&pkgLists.CreateParams{
				Title:   "Temp ListByWorkspace",
				BoardID: 2,
			})
			require.NoError(s.T(), err, "failed to create temp list")

			test.params.ID = tempList.ID
			list, err := s.uc.FullUpdate(test.params)
			assert.ErrorIs(s.T(), err, test.err, "unexpected error")

			if err == nil {
				// check updated list
				assert.Equal(s.T(), test.params.ID, list.ID, "incorrect ID")
				assert.Equal(s.T(), test.params.Title, list.Title, "incorrect Title")
				assert.Equal(s.T(), test.params.BoardID, list.BoardID, "incorrect BoardID")

				// check list in storages
				getList, err := s.uc.Get(list.ID)
				assert.NoError(s.T(), err, "failed to fetch list from the database")
				assert.Equal(s.T(), list.ID, getList.ID, "incorrect listID")
				assert.Equal(s.T(), list.BoardID, getList.BoardID, "incorrect BoardID")
				assert.Equal(s.T(), list.Title, getList.Title, "incorrect Title")
			}

			err = s.uc.Delete(tempList.ID)
			require.NoError(s.T(), err, "failed to delete temp list")
		})
	}
}

func (s *ListsSuite) TestPartialUpdate() {
	type testCase struct {
		params *pkgLists.PartialUpdateParams
		list   models.List
		err    error
	}

	tests := map[string]testCase{
		"full update": {
			params: &pkgLists.PartialUpdateParams{
				Title:       "MathStat",
				UpdateTitle: true,
			},
			list: models.List{
				Title:   "MathStat",
				BoardID: 2,
			},
			err: nil,
		},
		"only title update": {
			params: &pkgLists.PartialUpdateParams{
				Title:       "New MathStat",
				UpdateTitle: true,
			},
			list: models.List{
				Title:   "New MathStat",
				BoardID: 2,
			},
			err: nil,
		},
	}

	for name, test := range tests {
		s.Run(name, func() {
			tempList, err := s.uc.Create(&pkgLists.CreateParams{
				Title:   "Temp ListByWorkspace",
				BoardID: 2,
			})
			require.NoError(s.T(), err, "failed to create temp list")

			test.params.ID = tempList.ID
			list, err := s.uc.PartialUpdate(test.params)
			assert.ErrorIs(s.T(), err, test.err, "unexpected error")

			if err == nil {
				// check updated list
				assert.Equal(s.T(), test.params.ID, list.ID, "incorrect listID")
				assert.Equal(s.T(), test.list.Title, list.Title, "incorrect Title")
				assert.Equal(s.T(), test.list.BoardID, list.BoardID, "incorrect BoardID")

				// check list in storages
				getList, err := s.uc.Get(list.ID)
				assert.NoError(s.T(), err, "failed to fetch list from the database")
				assert.Equal(s.T(), test.list.Title, getList.Title, "incorrect Title")
				assert.Equal(s.T(), test.list.BoardID, getList.BoardID, "incorrect BoardID")
			}

			err = s.uc.Delete(tempList.ID)
			require.NoError(s.T(), err, "failed to delete temp list")
		})
	}
}

func (s *ListsSuite) TestDelete() {
	type testCase struct {
		setupList func() (models.List, error)
		err       error
	}

	tests := map[string]testCase{
		"normal": {
			setupList: func() (models.List, error) {
				return s.uc.Create(&pkgLists.CreateParams{
					Title:   "Test ListByWorkspace",
					BoardID: 1,
				})
			},
			err: nil,
		},
		"list not found": {
			setupList: func() (models.List, error) {
				return models.List{ID: 999}, nil
			},
			err: pkgErrors.ErrListNotFound,
		},
	}

	for name, test := range tests {
		s.Run(name, func() {
			list, err := test.setupList()
			s.Require().NoError(err)

			err = s.uc.Delete(list.ID)
			assert.ErrorIs(s.T(), err, test.err, "unexpected error")

			if test.err == nil {
				_, err = s.uc.Get(list.ID)
				assert.ErrorIs(s.T(), err, pkgErrors.ErrListNotFound, "list should be deleted")
			}
		})
	}
}

func TestListSuite(t *testing.T) {
	suite.Run(t, new(ListsSuite))
}
