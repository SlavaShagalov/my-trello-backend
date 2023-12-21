package boards

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/pkg/errors"

	pkgBoards "github.com/SlavaShagalov/my-trello-backend/internal/boards"
	"github.com/SlavaShagalov/my-trello-backend/internal/boards/mocks"
	boardsUsecase "github.com/SlavaShagalov/my-trello-backend/internal/boards/usecase"
	imgMocks "github.com/SlavaShagalov/my-trello-backend/internal/images/mocks"
	"github.com/SlavaShagalov/my-trello-backend/internal/models"
	pkgErrors "github.com/SlavaShagalov/my-trello-backend/internal/pkg/errors"
)

type BoardsUsecaseSuite struct {
	suite.Suite
}

func (s *BoardsUsecaseSuite) TestCreate(t provider.T) {
	type fields struct {
		repo    *mocks.MockRepository
		imgRepo *imgMocks.MockRepository
		params  *pkgBoards.CreateParams
		board   *models.Board
	}

	type testCase struct {
		prepare func(f *fields)
		params  *pkgBoards.CreateParams
		board   models.Board
		err     error
	}

	tests := map[string]testCase{
		"normal": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Create(f.params).Return(*f.board, nil)
			},
			params: &pkgBoards.CreateParams{
				Title:       "University",
				Description: "University Board",
				WorkspaceID: 27,
			},
			board: models.Board{
				ID:          210,
				WorkspaceID: 27,
				Title:       "University",
				Description: "University Board",
			},
			err: nil,
		},
		"workspace not found": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Create(f.params).Return(*f.board, pkgErrors.ErrWorkspaceNotFound)
			},
			params: &pkgBoards.CreateParams{Title: "University", Description: "University Board", WorkspaceID: 27},
			board:  models.Board{},
			err:    pkgErrors.ErrWorkspaceNotFound,
		},
		"storages error": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Create(f.params).Return(*f.board, pkgErrors.ErrDb)
			},
			params: &pkgBoards.CreateParams{Title: "University", Description: "University Board", WorkspaceID: 27},
			board:  models.Board{},
			err:    pkgErrors.ErrDb,
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
				params:  test.params, board: &test.board,
			}
			if test.prepare != nil {
				test.prepare(&f)
			}

			uc := boardsUsecase.New(f.repo, f.imgRepo)
			board, err := uc.Create(test.params)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if board != test.board {
				t.Errorf("\nExpected: %v\nGot: %v", test.board, board)
			}
		})
	}
}

func (s *BoardsUsecaseSuite) TestList(t provider.T) {
	type fields struct {
		repo        *mocks.MockRepository
		imgRepo     *imgMocks.MockRepository
		workspaceID int
		boards      []models.Board
	}

	type testCase struct {
		prepare     func(f *fields)
		workspaceID int
		boards      []models.Board
		err         error
	}

	tests := map[string]testCase{
		"normal": {
			prepare: func(f *fields) {
				f.repo.EXPECT().List(f.workspaceID).Return(f.boards, nil)
			},
			workspaceID: 27,
			boards: []models.Board{
				{ID: 21, WorkspaceID: 27, Title: "University", Description: "BMSTU Board"},
				{ID: 22, WorkspaceID: 27, Title: "Life", Description: "Life Board"},
				{ID: 23, WorkspaceID: 27, Title: "Work", Description: "Work Board"},
			},
			err: nil,
		},
		"empty result": {
			prepare: func(f *fields) {
				f.repo.EXPECT().List(f.workspaceID).Return(f.boards, nil)
			},
			workspaceID: 27,
			boards:      []models.Board{},
			err:         nil,
		},
		"board not found": {
			prepare: func(f *fields) {
				f.repo.EXPECT().List(f.workspaceID).Return(f.boards, pkgErrors.ErrWorkspaceNotFound)
			},
			workspaceID: 27,
			boards:      nil,
			err:         pkgErrors.ErrWorkspaceNotFound,
		},
		"storages error": {
			prepare: func(f *fields) {
				f.repo.EXPECT().List(f.workspaceID).Return(f.boards, pkgErrors.ErrDb)
			},
			workspaceID: 27,
			boards:      nil,
			err:         pkgErrors.ErrDb,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t provider.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				repo:        mocks.NewMockRepository(ctrl),
				imgRepo:     imgMocks.NewMockRepository(ctrl),
				workspaceID: test.workspaceID, boards: test.boards}
			if test.prepare != nil {
				test.prepare(&f)
			}

			serv := boardsUsecase.New(f.repo, f.imgRepo)
			boards, err := serv.ListByWorkspace(test.workspaceID)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if !reflect.DeepEqual(boards, test.boards) {
				t.Errorf("\nExpected: %v\nGot: %v", test.boards, boards)
			}
		})
	}
}

func (s *BoardsUsecaseSuite) TestGet(t provider.T) {
	type fields struct {
		repo    *mocks.MockRepository
		imgRepo *imgMocks.MockRepository
		id      int
		board   *models.Board
	}

	type testCase struct {
		prepare func(f *fields)
		id      int
		board   models.Board
		err     error
	}

	tests := map[string]testCase{
		"normal": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Get(f.id).Return(*f.board, nil)
			},
			id: 21,
			board: models.Board{
				ID:          21,
				WorkspaceID: 27,
				Title:       "University",
				Description: "University Board",
			},
			err: nil,
		},
		"board not found": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Get(f.id).Return(*f.board, pkgErrors.ErrBoardNotFound)
			},
			id:    21,
			board: models.Board{},
			err:   pkgErrors.ErrBoardNotFound,
		},
		"storages error": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Get(f.id).Return(*f.board, pkgErrors.ErrDb)
			},
			id:    21,
			board: models.Board{},
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
				id:      test.id, board: &test.board}
			if test.prepare != nil {
				test.prepare(&f)
			}

			uc := boardsUsecase.New(f.repo, f.imgRepo)
			board, err := uc.Get(test.id)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if board != test.board {
				t.Errorf("\nExpected: %v\nGot: %v", test.board, board)
			}
		})
	}
}

func (s *BoardsUsecaseSuite) TestFullUpdate(t provider.T) {
	type fields struct {
		repo    *mocks.MockRepository
		imgRepo *imgMocks.MockRepository
		params  *pkgBoards.FullUpdateParams
		board   *models.Board
	}

	type testCase struct {
		prepare func(f *fields)
		params  *pkgBoards.FullUpdateParams
		board   models.Board
		err     error
	}

	tests := map[string]testCase{
		"normal": {
			prepare: func(f *fields) {
				f.repo.EXPECT().FullUpdate(f.params).Return(*f.board, nil)
			},
			params: &pkgBoards.FullUpdateParams{
				ID:          21,
				Title:       "University",
				Description: "University Board",
				WorkspaceID: 27,
			},
			board: models.Board{ID: 21, WorkspaceID: 27, Title: "University", Description: "University Board"},
			err:   nil,
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
				params:  test.params, board: &test.board}
			if test.prepare != nil {
				test.prepare(&f)
			}

			uc := boardsUsecase.New(f.repo, f.imgRepo)
			board, err := uc.FullUpdate(test.params)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if board != test.board {
				t.Errorf("\nExpected: %v\nGot: %v", test.board, board)
			}
		})
	}
}

func (s *BoardsUsecaseSuite) TestPartialUpdate(t provider.T) {
	type fields struct {
		repo    *mocks.MockRepository
		imgRepo *imgMocks.MockRepository
		params  *pkgBoards.PartialUpdateParams
		board   *models.Board
	}

	type testCase struct {
		prepare func(f *fields)
		params  *pkgBoards.PartialUpdateParams
		board   models.Board
		err     error
	}

	tests := map[string]testCase{
		"normal": {
			prepare: func(f *fields) {
				f.repo.EXPECT().PartialUpdate(f.params).Return(*f.board, nil)
			},
			params: &pkgBoards.PartialUpdateParams{
				ID:                21,
				Title:             "University",
				UpdateTitle:       true,
				Description:       "University Board",
				UpdateDescription: true,
				WorkspaceID:       27,
				UpdateWorkspaceID: true,
			},
			board: models.Board{ID: 21, WorkspaceID: 27, Title: "University", Description: "University Board"},
			err:   nil,
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
				params:  test.params, board: &test.board}
			if test.prepare != nil {
				test.prepare(&f)
			}

			uc := boardsUsecase.New(f.repo, f.imgRepo)
			board, err := uc.PartialUpdate(test.params)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if board != test.board {
				t.Errorf("\nExpected: %v\nGot: %v", test.board, board)
			}
		})
	}
}

func (s *BoardsUsecaseSuite) TestDelete(t provider.T) {
	type fields struct {
		repo    *mocks.MockRepository
		imgRepo *imgMocks.MockRepository
		id      int
	}

	type testCase struct {
		prepare func(f *fields)
		id      int
		err     error
	}

	tests := map[string]testCase{
		"normal": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Delete(f.id).Return(nil)
			},
			id:  21,
			err: nil,
		},
		"board not found": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Delete(f.id).Return(pkgErrors.ErrBoardNotFound)
			},
			id:  21,
			err: pkgErrors.ErrBoardNotFound,
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
				id:      test.id}
			if test.prepare != nil {
				test.prepare(&f)
			}

			uc := boardsUsecase.New(f.repo, f.imgRepo)
			err := uc.Delete(test.id)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
		})
	}
}

func TestSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(BoardsUsecaseSuite))
}
