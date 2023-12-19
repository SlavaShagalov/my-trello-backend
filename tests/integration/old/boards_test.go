package old

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

	pkgBoards "github.com/SlavaShagalov/my-trello-backend/internal/boards"
	pkgZap "github.com/SlavaShagalov/my-trello-backend/internal/pkg/log/zap"
	pkgDb "github.com/SlavaShagalov/my-trello-backend/internal/pkg/storages"

	boardsRepo "github.com/SlavaShagalov/my-trello-backend/internal/boards/repository/postgres"
	boardsUC "github.com/SlavaShagalov/my-trello-backend/internal/boards/usecase"
)

type BoardsSuite struct {
	suite.Suite
	db     *sql.DB
	logger *zap.Logger
	uc     pkgBoards.Usecase
}

func (s *BoardsSuite) SetupSuite() {
	s.logger = pkgZap.NewTestLogger()

	var err error
	config.SetTestPostgresConfig()
	s.db, err = pkgDb.NewPostgres(s.logger)
	s.Require().NoError(err)

	repo := boardsRepo.New(s.db, s.logger)
	s.uc = boardsUC.New(repo)
}

func (s *BoardsSuite) TearDownSuite() {
	err := s.db.Close()
	s.Require().NoError(err)

	_ = s.logger.Sync()
}

func (s *BoardsSuite) TestCreate() {
	type testCase struct {
		params *pkgBoards.CreateParams
		err    error
	}

	tests := map[string]testCase{
		"normal": {
			params: &pkgBoards.CreateParams{
				Title:       "University",
				Description: "BMSTU board",
				WorkspaceID: 3,
			},
			err: nil,
		},
		"workspace not found": {
			params: &pkgBoards.CreateParams{
				Title:       "University",
				Description: "BMSTU board",
				WorkspaceID: 999,
			},
			err: pkgErrors.ErrWorkspaceNotFound,
		},
	}

	for name, test := range tests {
		s.Run(name, func() {
			board, err := s.uc.Create(test.params)
			assert.ErrorIs(s.T(), err, test.err, "unexpected error")

			if err == nil {
				assert.Equal(s.T(), test.params.WorkspaceID, board.WorkspaceID, "incorrect WorkspaceID")
				assert.Equal(s.T(), test.params.Title, board.Title, "incorrect Title")
				assert.Equal(s.T(), test.params.Description, board.Description, "incorrect Description")

				getBoard, err := s.uc.Get(board.ID)
				assert.NoError(s.T(), err, "failed to fetch board from the database")
				assert.Equal(s.T(), board.ID, getBoard.ID, "incorrect boardID")
				assert.Equal(s.T(), test.params.WorkspaceID, getBoard.WorkspaceID, "incorrect WorkspaceID")
				assert.Equal(s.T(), test.params.Title, getBoard.Title, "incorrect Title")
				assert.Equal(s.T(), test.params.Description, getBoard.Description, "incorrect Description")

				err = s.uc.Delete(board.ID)
				assert.NoError(s.T(), err, "failed to delete created board")
			}
		})
	}
}

func (s *BoardsSuite) TestList() {
	type testCase struct {
		userID int
		boards []models.Board
		err    error
	}

	tests := map[string]testCase{
		"normal": {
			userID: 2,
			boards: []models.Board{
				{
					ID:          4,
					WorkspaceID: 2,
					Title:       "Маркетинг и продвижение",
					Description: "Доска для планирования маркетинговых мероприятий проекта \"Бета\"",
				},
				{
					ID:          5,
					WorkspaceID: 2,
					Title:       "Анализ рынка",
					Description: "Доска для анализа рынка и конкурентов проекта \"Бета\"",
				},
				{
					ID:          6,
					WorkspaceID: 2,
					Title:       "Отчетность и аналитика",
					Description: "Доска для отчетности и анализа результатов проекта \"Бета\"",
				},
			},
			err: nil,
		},
		"empty result": {
			userID: 8,
			boards: []models.Board{},
			err:    nil,
		},
	}

	for name, test := range tests {
		s.Run(name, func() {
			boards, err := s.uc.ListByWorkspace(test.userID)

			assert.ErrorIs(s.T(), err, test.err, "unexpected error")

			if err == nil {
				assert.Equal(s.T(), len(test.boards), len(boards), "incorrect boards length")
				for i := 0; i < len(test.boards); i++ {
					assert.Equal(s.T(), test.boards[i].ID, boards[i].ID, "incorrect boardID")
					assert.Equal(s.T(), test.boards[i].WorkspaceID, boards[i].WorkspaceID, "incorrect WorkspaceID")
					assert.Equal(s.T(), test.boards[i].Title, boards[i].Title, "incorrect Title")
					assert.Equal(s.T(), test.boards[i].Description, boards[i].Description, "incorrect Description")
				}
			}
		})
	}
}

func (s *BoardsSuite) TestGet() {
	type testCase struct {
		boardID int
		board   models.Board
		err     error
	}

	tests := map[string]testCase{
		"normal": {
			boardID: 8,
			board: models.Board{
				ID:          8,
				WorkspaceID: 3,
				Title:       "Исследование пользователей",
				Description: "Доска для проведения исследований пользователей проекта \"Гамма\"",
			},
			err: nil,
		},
		"board not found": {
			boardID: 999,
			board:   models.Board{},
			err:     pkgErrors.ErrBoardNotFound,
		},
	}

	for name, test := range tests {
		s.Run(name, func() {
			board, err := s.uc.Get(test.boardID)

			assert.ErrorIs(s.T(), err, test.err, "unexpected error")

			if err == nil {
				assert.Equal(s.T(), test.board.ID, board.ID, "incorrect boardID")
				assert.Equal(s.T(), test.board.WorkspaceID, board.WorkspaceID, "incorrect WorkspaceID")
				assert.Equal(s.T(), test.board.Title, board.Title, "incorrect Title")
				assert.Equal(s.T(), test.board.Description, board.Description, "incorrect Description")
			}
		})
	}
}

func (s *BoardsSuite) TestFullUpdate() {
	type testCase struct {
		params *pkgBoards.FullUpdateParams
		err    error
	}

	tests := map[string]testCase{
		"normal": {
			params: &pkgBoards.FullUpdateParams{
				Title:       "University",
				Description: "BMSTU board",
				WorkspaceID: 3,
			},
			err: nil,
		},
	}

	for name, test := range tests {
		s.Run(name, func() {
			tempBoard, err := s.uc.Create(&pkgBoards.CreateParams{
				Title:       "Temp Board",
				Description: "Temp Board Description",
				WorkspaceID: 2,
			})
			require.NoError(s.T(), err, "failed to create temp board")

			test.params.ID = tempBoard.ID
			board, err := s.uc.FullUpdate(test.params)
			assert.ErrorIs(s.T(), err, test.err, "unexpected error")

			if err == nil {
				// check updated board
				assert.Equal(s.T(), test.params.ID, board.ID, "incorrect ID")
				assert.Equal(s.T(), test.params.Title, board.Title, "incorrect Title")
				assert.Equal(s.T(), test.params.Description, board.Description, "incorrect Description")
				assert.Equal(s.T(), test.params.WorkspaceID, board.WorkspaceID, "incorrect WorkspaceID")

				// check board in storages
				getBoard, err := s.uc.Get(board.ID)
				assert.NoError(s.T(), err, "failed to fetch board from the database")
				assert.Equal(s.T(), board.ID, getBoard.ID, "incorrect boardID")
				assert.Equal(s.T(), board.WorkspaceID, getBoard.WorkspaceID, "incorrect WorkspaceID")
				assert.Equal(s.T(), board.Title, getBoard.Title, "incorrect Title")
				assert.Equal(s.T(), board.Description, getBoard.Description, "incorrect Description")
			}

			err = s.uc.Delete(tempBoard.ID)
			require.NoError(s.T(), err, "failed to delete temp board")
		})
	}
}

func (s *BoardsSuite) TestPartialUpdate() {
	type testCase struct {
		params *pkgBoards.PartialUpdateParams
		board  models.Board
		err    error
	}

	tests := map[string]testCase{
		"full update": {
			params: &pkgBoards.PartialUpdateParams{
				Title:             "University",
				UpdateTitle:       true,
				Description:       "BMSTU board",
				UpdateDescription: true,
			},
			board: models.Board{
				Title:       "University",
				Description: "BMSTU board",
				WorkspaceID: 2,
			},
			err: nil,
		},
		"only title update": {
			params: &pkgBoards.PartialUpdateParams{
				Title:       "New University",
				UpdateTitle: true,
			},
			board: models.Board{
				Title:       "New University",
				Description: "Temp Board Description",
				WorkspaceID: 2,
			},
			err: nil,
		},
		"only description update": {
			params: &pkgBoards.PartialUpdateParams{
				Description:       "New BMSTU board",
				UpdateDescription: true,
			},
			board: models.Board{
				Title:       "Temp Board",
				Description: "New BMSTU board",
				WorkspaceID: 2,
			},
			err: nil,
		},
	}

	for name, test := range tests {
		s.Run(name, func() {
			tempBoard, err := s.uc.Create(&pkgBoards.CreateParams{
				Title:       "Temp Board",
				Description: "Temp Board Description",
				WorkspaceID: 2,
			})
			require.NoError(s.T(), err, "failed to create temp board")

			test.params.ID = tempBoard.ID
			board, err := s.uc.PartialUpdate(test.params)
			assert.ErrorIs(s.T(), err, test.err, "unexpected error")

			if err == nil {
				// check updated board
				assert.Equal(s.T(), test.params.ID, board.ID, "incorrect boardID")
				assert.Equal(s.T(), test.board.Title, board.Title, "incorrect Title")
				assert.Equal(s.T(), test.board.Description, board.Description, "incorrect Description")
				assert.Equal(s.T(), test.board.WorkspaceID, board.WorkspaceID, "incorrect WorkspaceID")

				// check board in storages
				getBoard, err := s.uc.Get(board.ID)
				assert.NoError(s.T(), err, "failed to fetch board from the database")
				assert.Equal(s.T(), test.board.Title, getBoard.Title, "incorrect Title")
				assert.Equal(s.T(), test.board.Description, getBoard.Description, "incorrect Description")
				assert.Equal(s.T(), test.board.WorkspaceID, getBoard.WorkspaceID, "incorrect WorkspaceID")
			}

			err = s.uc.Delete(tempBoard.ID)
			require.NoError(s.T(), err, "failed to delete temp board")
		})
	}
}

func (s *BoardsSuite) TestDelete() {
	type testCase struct {
		setupBoard func() (models.Board, error)
		err        error
	}

	tests := map[string]testCase{
		"normal": {
			setupBoard: func() (models.Board, error) {
				return s.uc.Create(&pkgBoards.CreateParams{
					Title:       "Test Board",
					Description: "Test Board Description",
					WorkspaceID: 1,
				})
			},
			err: nil,
		},
		"board not found": {
			setupBoard: func() (models.Board, error) {
				return models.Board{ID: 999}, nil
			},
			err: pkgErrors.ErrBoardNotFound,
		},
	}

	for name, test := range tests {
		s.Run(name, func() {
			board, err := test.setupBoard()
			s.Require().NoError(err)

			err = s.uc.Delete(board.ID)
			assert.ErrorIs(s.T(), err, test.err, "unexpected error")

			if test.err == nil {
				_, err = s.uc.Get(board.ID)
				assert.ErrorIs(s.T(), err, pkgErrors.ErrBoardNotFound, "board should be deleted")
			}
		})
	}
}

func TestBoardSuite(t *testing.T) {
	suite.Run(t, new(BoardsSuite))
}
