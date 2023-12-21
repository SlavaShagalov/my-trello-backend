package lists

import (
	pkgLists "github.com/SlavaShagalov/my-trello-backend/internal/lists"
	"github.com/SlavaShagalov/my-trello-backend/internal/lists/mocks"
	listsUsecase "github.com/SlavaShagalov/my-trello-backend/internal/lists/usecase"
	"github.com/SlavaShagalov/my-trello-backend/internal/models"
	pkgErrors "github.com/SlavaShagalov/my-trello-backend/internal/pkg/errors"
	"github.com/SlavaShagalov/my-trello-backend/tests/utils/builder"
	"github.com/golang/mock/gomock"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/pkg/errors"
	"reflect"
	"testing"
)

type ListsUsecaseSuite struct {
	suite.Suite

	listsBuilder *builder.ListBuilder
}

func (s *ListsUsecaseSuite) BeforeEach(t provider.T) {
	t.WithNewStep("SetupTest step", func(ctx provider.StepCtx) {})

	s.listsBuilder = builder.NewListBuilder()
}

func (s *ListsUsecaseSuite) TestCreate(t provider.T) {
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
			list: s.listsBuilder.
				WithID(21).
				WithBoardID(27).
				WithTitle("MathStat").
				WithPosition(41).
				Build(),
			err: nil,
		},
		"board not found": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Create(f.params).Return(*f.list, pkgErrors.ErrBoardNotFound)
			},
			params: &pkgLists.CreateParams{Title: "MathStat", BoardID: 27},
			list:   s.listsBuilder.Build(),
			err:    pkgErrors.ErrBoardNotFound,
		},
		"storages error": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Create(f.params).Return(*f.list, pkgErrors.ErrDb)
			},
			params: &pkgLists.CreateParams{Title: "MathStat", BoardID: 27},
			list:   s.listsBuilder.Build(),
			err:    pkgErrors.ErrDb,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t provider.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl), params: test.params, list: &test.list}
			if test.prepare != nil {
				test.prepare(&f)
			}

			uc := listsUsecase.New(f.repo)
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

func (s *ListsUsecaseSuite) TestList(t provider.T) {
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
				f.repo.EXPECT().ListByBoard(f.boardID).Return(f.lists, nil)
			},
			boardID: 27,
			lists: []models.List{
				s.listsBuilder.
					WithID(21).
					WithBoardID(27).
					WithTitle("MathStat").
					WithPosition(41).
					Build(),
				s.listsBuilder.
					WithID(22).
					WithBoardID(27).
					WithTitle("Software Design").
					WithPosition(42).
					Build(),
				s.listsBuilder.
					WithID(23).
					WithBoardID(27).
					WithTitle("Operating Systems").
					WithPosition(43).
					Build(),
			},
			err: nil,
		},
		"empty result": {
			prepare: func(f *fields) {
				f.repo.EXPECT().ListByBoard(f.boardID).Return(f.lists, nil)
			},
			boardID: 27,
			lists:   []models.List{},
			err:     nil,
		},
		"board not found": {
			prepare: func(f *fields) {
				f.repo.EXPECT().ListByBoard(f.boardID).Return(f.lists, pkgErrors.ErrBoardNotFound)
			},
			boardID: 27,
			lists:   nil,
			err:     pkgErrors.ErrBoardNotFound,
		},
		"storages error": {
			prepare: func(f *fields) {
				f.repo.EXPECT().ListByBoard(f.boardID).Return(f.lists, pkgErrors.ErrDb)
			},
			boardID: 27,
			lists:   nil,
			err:     pkgErrors.ErrDb,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t provider.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl), boardID: test.boardID, lists: test.lists}
			if test.prepare != nil {
				test.prepare(&f)
			}

			serv := listsUsecase.New(f.repo)
			lists, err := serv.ListByBoard(test.boardID)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if !reflect.DeepEqual(lists, test.lists) {
				t.Errorf("\nExpected: %v\nGot: %v", test.lists, lists)
			}
		})
	}
}

func (s *ListsUsecaseSuite) TestGet(t provider.T) {
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
			id: 21,
			list: s.listsBuilder.
				WithID(21).
				WithBoardID(27).
				WithTitle("MathStat").
				WithPosition(41).
				Build(),
			err: nil,
		},
		"list not found": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Get(f.id).Return(*f.list, pkgErrors.ErrListNotFound)
			},
			id:   21,
			list: s.listsBuilder.Build(),
			err:  pkgErrors.ErrListNotFound,
		},
		"db error": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Get(f.id).Return(*f.list, pkgErrors.ErrDb)
			},
			id:   21,
			list: s.listsBuilder.Build(),
			err:  pkgErrors.ErrDb,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t provider.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl), id: test.id, list: &test.list}
			if test.prepare != nil {
				test.prepare(&f)
			}

			uc := listsUsecase.New(f.repo)
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

func (s *ListsUsecaseSuite) TestFullUpdate(t provider.T) {
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
			list: s.listsBuilder.
				WithID(21).
				WithBoardID(27).
				WithTitle("MathStat").
				WithPosition(41).
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

			f := fields{repo: mocks.NewMockRepository(ctrl), params: test.params, list: &test.list}
			if test.prepare != nil {
				test.prepare(&f)
			}

			uc := listsUsecase.New(f.repo)
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

func (s *ListsUsecaseSuite) TestPartialUpdate(t provider.T) {
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
			list: s.listsBuilder.
				WithID(21).
				WithBoardID(27).
				WithTitle("MathStat").
				WithPosition(41).
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

			f := fields{repo: mocks.NewMockRepository(ctrl), params: test.params, list: &test.list}
			if test.prepare != nil {
				test.prepare(&f)
			}

			uc := listsUsecase.New(f.repo)
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

func (s *ListsUsecaseSuite) TestDelete(t provider.T) {
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
		t.Run(name, func(t provider.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl), id: test.id}
			if test.prepare != nil {
				test.prepare(&f)
			}

			uc := listsUsecase.New(f.repo)
			err := uc.Delete(test.id)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
		})
	}
}

func TestSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(ListsUsecaseSuite))
}
