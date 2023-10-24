package usecase

import (
	pkgBoards "github.com/SlavaShagalov/my-trello-backend/internal/boards"
	"github.com/SlavaShagalov/my-trello-backend/internal/boards/mocks"
	"github.com/SlavaShagalov/my-trello-backend/internal/models"
	pkgErrors "github.com/SlavaShagalov/my-trello-backend/internal/pkg/errors"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"reflect"
	"testing"
)

func TestUsecase_Create(t *testing.T) {
	type fields struct {
		repo   *mocks.MockRepository
		params *pkgBoards.CreateParams
		board  *models.Board
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
				ID:          21,
				WorkspaceID: 27,
				Title:       "University",
				Description: "University Board",
				Background:  "university.jpg",
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
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl), params: test.params, board: &test.board}
			if test.prepare != nil {
				test.prepare(&f)
			}

			uc := New(f.repo)
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

func TestUsecase_List(t *testing.T) {
	type fields struct {
		repo        *mocks.MockRepository
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
				{ID: 21, WorkspaceID: 27, Title: "University", Description: "BMSTU Board", Background: "university.jpg"},
				{ID: 22, WorkspaceID: 27, Title: "Life", Description: "Life Board", Background: "life.jpg"},
				{ID: 23, WorkspaceID: 27, Title: "Work", Description: "Work Board", Background: "work.jpg"},
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
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl), workspaceID: test.workspaceID, boards: test.boards}
			if test.prepare != nil {
				test.prepare(&f)
			}

			serv := New(f.repo)
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

func TestUsecase_Get(t *testing.T) {
	type fields struct {
		repo  *mocks.MockRepository
		id    int
		board *models.Board
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
				Background:  "university.jpg",
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
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl), id: test.id, board: &test.board}
			if test.prepare != nil {
				test.prepare(&f)
			}

			uc := New(f.repo)
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

func TestFullUpdate(t *testing.T) {
	type fields struct {
		repo   *mocks.MockRepository
		params *pkgBoards.FullUpdateParams
		board  *models.Board
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
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl), params: test.params, board: &test.board}
			if test.prepare != nil {
				test.prepare(&f)
			}

			uc := New(f.repo)
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

func TestPartialUpdate(t *testing.T) {
	type fields struct {
		repo   *mocks.MockRepository
		params *pkgBoards.PartialUpdateParams
		board  *models.Board
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
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl), params: test.params, board: &test.board}
			if test.prepare != nil {
				test.prepare(&f)
			}

			uc := New(f.repo)
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

func TestUsecase_Delete(t *testing.T) {
	type fields struct {
		repo *mocks.MockRepository
		id   int
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
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl), id: test.id}
			if test.prepare != nil {
				test.prepare(&f)
			}

			uc := New(f.repo)
			err := uc.Delete(test.id)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
		})
	}
}
