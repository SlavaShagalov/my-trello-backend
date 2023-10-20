package usecase

import (
	pkgLists "github.com/SlavaShagalov/my-trello-backend/internal/lists"
	"github.com/SlavaShagalov/my-trello-backend/internal/lists/mocks"
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
		params *pkgLists.CreateParams
		list   *models.List
	}

	type testCase struct {
		prepare func(f *fields)
		params  *pkgLists.CreateParams
		list    models.List
		err     error
	}

	tests := map[string]testCase{
		"normal": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Create(f.params).Return(*f.list, nil)
			},
			params: &pkgLists.CreateParams{
				Title:   "MathStat",
				BoardID: 27,
			},
			list: models.List{ID: 21, BoardID: 27, Title: "MathStat", Position: 41},
			err:  nil,
		},
		"board not found": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Create(f.params).Return(*f.list, pkgErrors.ErrBoardNotFound)
			},
			params: &pkgLists.CreateParams{Title: "MathStat", BoardID: 27},
			list:   models.List{},
			err:    pkgErrors.ErrBoardNotFound,
		},
		"storages error": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Create(f.params).Return(*f.list, pkgErrors.ErrDb)
			},
			params: &pkgLists.CreateParams{Title: "MathStat", BoardID: 27},
			list:   models.List{},
			err:    pkgErrors.ErrDb,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl), params: test.params, list: &test.list}
			if test.prepare != nil {
				test.prepare(&f)
			}

			uc := NewUsecase(f.repo)
			list, err := uc.Create(test.params)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if list != test.list {
				t.Errorf("\nExpected: %v\nGot: %v", test.list, list)
			}
		})
	}
}

func TestUsecase_List(t *testing.T) {
	type fields struct {
		repo    *mocks.MockRepository
		boardID int
		lists   []models.List
	}

	type testCase struct {
		prepare func(f *fields)
		boardID int
		lists   []models.List
		err     error
	}

	tests := map[string]testCase{
		"normal": {
			prepare: func(f *fields) {
				f.repo.EXPECT().List(f.boardID).Return(f.lists, nil)
			},
			boardID: 27,
			lists: []models.List{
				{ID: 21, BoardID: 27, Title: "MathStat", Position: 41},
				{ID: 22, BoardID: 27, Title: "Software Design", Position: 42},
				{ID: 23, BoardID: 27, Title: "Operating Systems", Position: 43},
			},
			err: nil,
		},
		"empty result": {
			prepare: func(f *fields) {
				f.repo.EXPECT().List(f.boardID).Return(f.lists, nil)
			},
			boardID: 27,
			lists:   []models.List{},
			err:     nil,
		},
		"board not found": {
			prepare: func(f *fields) {
				f.repo.EXPECT().List(f.boardID).Return(f.lists, pkgErrors.ErrBoardNotFound)
			},
			boardID: 27,
			lists:   nil,
			err:     pkgErrors.ErrBoardNotFound,
		},
		"storages error": {
			prepare: func(f *fields) {
				f.repo.EXPECT().List(f.boardID).Return(f.lists, pkgErrors.ErrDb)
			},
			boardID: 27,
			lists:   nil,
			err:     pkgErrors.ErrDb,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl), boardID: test.boardID, lists: test.lists}
			if test.prepare != nil {
				test.prepare(&f)
			}

			serv := NewUsecase(f.repo)
			lists, err := serv.List(test.boardID)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if !reflect.DeepEqual(lists, test.lists) {
				t.Errorf("\nExpected: %v\nGot: %v", test.lists, lists)
			}
		})
	}
}

func TestUsecase_Get(t *testing.T) {
	type fields struct {
		repo *mocks.MockRepository
		id   int
		list *models.List
	}

	type testCase struct {
		prepare func(f *fields)
		id      int
		list    models.List
		err     error
	}

	tests := map[string]testCase{
		"normal": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Get(f.id).Return(*f.list, nil)
			},
			id:   21,
			list: models.List{ID: 21, BoardID: 27, Title: "MathStat", Position: 41},
			err:  nil,
		},
		"list not found": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Get(f.id).Return(*f.list, pkgErrors.ErrListNotFound)
			},
			id:   21,
			list: models.List{},
			err:  pkgErrors.ErrListNotFound,
		},
		"storages error": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Get(f.id).Return(*f.list, pkgErrors.ErrDb)
			},
			id:   21,
			list: models.List{},
			err:  pkgErrors.ErrDb,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl), id: test.id, list: &test.list}
			if test.prepare != nil {
				test.prepare(&f)
			}

			uc := NewUsecase(f.repo)
			list, err := uc.Get(test.id)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if list != test.list {
				t.Errorf("\nExpected: %v\nGot: %v", test.list, list)
			}
		})
	}
}

func TestFullUpdate(t *testing.T) {
	type fields struct {
		repo   *mocks.MockRepository
		params *pkgLists.FullUpdateParams
		list   *models.List
	}

	type testCase struct {
		prepare func(f *fields)
		params  *pkgLists.FullUpdateParams
		list    models.List
		err     error
	}

	tests := map[string]testCase{
		"normal": {
			prepare: func(f *fields) {
				f.repo.EXPECT().FullUpdate(f.params).Return(*f.list, nil)
			},
			params: &pkgLists.FullUpdateParams{ID: 21, Title: "MathStat", Position: 41, BoardID: 27},
			list:   models.List{ID: 21, BoardID: 27, Title: "MathStat", Position: 41},
			err:    nil,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl), params: test.params, list: &test.list}
			if test.prepare != nil {
				test.prepare(&f)
			}

			uc := NewUsecase(f.repo)
			list, err := uc.FullUpdate(test.params)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if list != test.list {
				t.Errorf("\nExpected: %v\nGot: %v", test.list, list)
			}
		})
	}
}

func TestPartialUpdate(t *testing.T) {
	type fields struct {
		repo   *mocks.MockRepository
		params *pkgLists.PartialUpdateParams
		list   *models.List
	}

	type testCase struct {
		prepare func(f *fields)
		params  *pkgLists.PartialUpdateParams
		list    models.List
		err     error
	}

	tests := map[string]testCase{
		"normal": {
			prepare: func(f *fields) {
				f.repo.EXPECT().PartialUpdate(f.params).Return(*f.list, nil)
			},
			params: &pkgLists.PartialUpdateParams{
				ID:             21,
				Title:          "MathStat",
				UpdateTitle:    true,
				Position:       41,
				UpdatePosition: true,
				BoardID:        27,
				UpdateBoardID:  true,
			},
			list: models.List{ID: 21, BoardID: 27, Title: "MathStat", Position: 41},
			err:  nil,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl), params: test.params, list: &test.list}
			if test.prepare != nil {
				test.prepare(&f)
			}

			uc := NewUsecase(f.repo)
			list, err := uc.PartialUpdate(test.params)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if list != test.list {
				t.Errorf("\nExpected: %v\nGot: %v", test.list, list)
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
		"list not found": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Delete(f.id).Return(pkgErrors.ErrListNotFound)
			},
			id:  21,
			err: pkgErrors.ErrListNotFound,
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

			uc := NewUsecase(f.repo)
			err := uc.Delete(test.id)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
		})
	}
}
