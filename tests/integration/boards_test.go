package integration

import (
	"context"
	"database/sql"
	pkgBoards "github.com/SlavaShagalov/my-trello-backend/internal/boards"
	"github.com/SlavaShagalov/my-trello-backend/internal/models"
	"github.com/SlavaShagalov/my-trello-backend/internal/pkg/config"
	pkgErrors "github.com/SlavaShagalov/my-trello-backend/internal/pkg/errors"
	pkgZap "github.com/SlavaShagalov/my-trello-backend/internal/pkg/log/zap"
	"github.com/SlavaShagalov/my-trello-backend/internal/pkg/opentel"
	pkgDb "github.com/SlavaShagalov/my-trello-backend/internal/pkg/storages/postgres"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
	"log"
	"os"
	"testing"

	boardsRepo "github.com/SlavaShagalov/my-trello-backend/internal/boards/repository/std"
	boardsUC "github.com/SlavaShagalov/my-trello-backend/internal/boards/usecase"
	imgMocks "github.com/SlavaShagalov/my-trello-backend/internal/images/mocks"
)

type BoardsSuite struct {
	suite.Suite
	db      *sql.DB
	log     *zap.Logger
	logfile *os.File
	uc      pkgBoards.Usecase
	tp      *sdktrace.TracerProvider
	mp      *sdkmetric.MeterProvider
	ctx     context.Context
}

func (s *BoardsSuite) SetupSuite() {
	s.ctx = context.Background()

	var err error
	s.log, s.logfile, err = pkgZap.NewTestLogger("/logs/boards.log")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	config.SetTestPostgresConfig()
	s.db, err = pkgDb.NewStd(s.log)
	s.Require().NoError(err)

	// Set up OpenTelemetry.
	serviceName := "test"
	serviceVersion := "0.1.0"
	s.tp, s.mp, err = opentel.SetupOTelSDK(context.Background(), s.log, serviceName, serviceVersion)
	if err != nil {
		return
	}

	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	repo := boardsRepo.New(s.db, s.log)
	imgRepo := imgMocks.NewMockRepository(ctrl)
	s.uc = boardsUC.New(repo, imgRepo)
}

func (s *BoardsSuite) TearDownSuite() {
	err := s.db.Close()
	s.Require().NoError(err)
	s.log.Info("DB connection closed")

	err = s.log.Sync()
	if err != nil {
		log.Println(err)
	}
	err = s.logfile.Close()
	if err != nil {
		log.Println(err)
	}

	_ = s.tp.Shutdown(s.ctx)
	_ = s.mp.Shutdown(s.ctx)
	s.log.Info("OpenTelemetry shutdown")
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
			ctx, span := opentel.Tracer.Start(context.Background(), "TestCreate "+name)
			defer span.End()

			board, err := s.uc.Create(ctx, test.params)
			assert.ErrorIs(s.T(), err, test.err, "unexpected error")

			if err == nil {
				assert.Equal(s.T(), test.params.WorkspaceID, board.WorkspaceID, "incorrect WorkspaceID")
				assert.Equal(s.T(), test.params.Title, board.Title, "incorrect Title")
				assert.Equal(s.T(), test.params.Description, board.Description, "incorrect Description")

				getBoard, err := s.uc.Get(ctx, board.ID)
				assert.NoError(s.T(), err, "failed to fetch board from the database")
				assert.Equal(s.T(), board.ID, getBoard.ID, "incorrect boardID")
				assert.Equal(s.T(), test.params.WorkspaceID, getBoard.WorkspaceID, "incorrect WorkspaceID")
				assert.Equal(s.T(), test.params.Title, getBoard.Title, "incorrect Title")
				assert.Equal(s.T(), test.params.Description, getBoard.Description, "incorrect Description")

				err = s.uc.Delete(ctx, board.ID)
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
			ctx, span := opentel.Tracer.Start(context.Background(), "TestList "+name)
			defer span.End()

			boards, err := s.uc.ListByWorkspace(ctx, test.userID)

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
			ctx, span := opentel.Tracer.Start(context.Background(), "TestGet "+name)
			defer span.End()

			board, err := s.uc.Get(ctx, test.boardID)

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
			ctx, span := opentel.Tracer.Start(context.Background(), "TestFullUpdate "+name)
			defer span.End()

			tempBoard, err := s.uc.Create(ctx, &pkgBoards.CreateParams{
				Title:       "Temp Board",
				Description: "Temp Board Description",
				WorkspaceID: 2,
			})
			require.NoError(s.T(), err, "failed to create temp board")

			test.params.ID = tempBoard.ID
			board, err := s.uc.FullUpdate(ctx, test.params)
			assert.ErrorIs(s.T(), err, test.err, "unexpected error")

			if err == nil {
				// check updated board
				assert.Equal(s.T(), test.params.ID, board.ID, "incorrect ID")
				assert.Equal(s.T(), test.params.Title, board.Title, "incorrect Title")
				assert.Equal(s.T(), test.params.Description, board.Description, "incorrect Description")
				assert.Equal(s.T(), test.params.WorkspaceID, board.WorkspaceID, "incorrect WorkspaceID")

				// check board in storages
				getBoard, err := s.uc.Get(ctx, board.ID)
				assert.NoError(s.T(), err, "failed to fetch board from the database")
				assert.Equal(s.T(), board.ID, getBoard.ID, "incorrect boardID")
				assert.Equal(s.T(), board.WorkspaceID, getBoard.WorkspaceID, "incorrect WorkspaceID")
				assert.Equal(s.T(), board.Title, getBoard.Title, "incorrect Title")
				assert.Equal(s.T(), board.Description, getBoard.Description, "incorrect Description")
			}

			err = s.uc.Delete(ctx, tempBoard.ID)
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
			ctx, span := opentel.Tracer.Start(context.Background(), "TestPartialUpdate "+name)
			defer span.End()

			tempBoard, err := s.uc.Create(ctx, &pkgBoards.CreateParams{
				Title:       "Temp Board",
				Description: "Temp Board Description",
				WorkspaceID: 2,
			})
			require.NoError(s.T(), err, "failed to create temp board")

			test.params.ID = tempBoard.ID
			board, err := s.uc.PartialUpdate(ctx, test.params)
			assert.ErrorIs(s.T(), err, test.err, "unexpected error")

			if err == nil {
				// check updated board
				assert.Equal(s.T(), test.params.ID, board.ID, "incorrect boardID")
				assert.Equal(s.T(), test.board.Title, board.Title, "incorrect Title")
				assert.Equal(s.T(), test.board.Description, board.Description, "incorrect Description")
				assert.Equal(s.T(), test.board.WorkspaceID, board.WorkspaceID, "incorrect WorkspaceID")

				// check board in storages
				getBoard, err := s.uc.Get(ctx, board.ID)
				assert.NoError(s.T(), err, "failed to fetch board from the database")
				assert.Equal(s.T(), test.board.Title, getBoard.Title, "incorrect Title")
				assert.Equal(s.T(), test.board.Description, getBoard.Description, "incorrect Description")
				assert.Equal(s.T(), test.board.WorkspaceID, getBoard.WorkspaceID, "incorrect WorkspaceID")
			}

			err = s.uc.Delete(ctx, tempBoard.ID)
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
				return s.uc.Create(context.Background(), &pkgBoards.CreateParams{
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
			ctx, span := opentel.Tracer.Start(context.Background(), "TestDelete "+name)
			defer span.End()

			board, err := test.setupBoard()
			s.Require().NoError(err)

			err = s.uc.Delete(ctx, board.ID)
			assert.ErrorIs(s.T(), err, test.err, "unexpected error")

			if test.err == nil {
				_, err = s.uc.Get(ctx, board.ID)
				assert.ErrorIs(s.T(), err, pkgErrors.ErrBoardNotFound, "board should be deleted")
			}
		})
	}
}

func TestBoardSuite(t *testing.T) {
	suite.Run(t, new(BoardsSuite))
}
