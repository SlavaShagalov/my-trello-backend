package usecase

import (
	pkgCards "github.com/SlavaShagalov/my-trello-backend/internal/cards"
	"github.com/SlavaShagalov/my-trello-backend/internal/cards/mocks"
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
		params *pkgCards.CreateParams
		card   *models.Card
	}

	type testCase struct {
		prepare func(f *fields)
		params  *pkgCards.CreateParams
		card    models.Card
		err     error
	}

	tests := map[string]testCase{
		"normal": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Create(f.params).Return(*f.card, nil)
			},
			params: &pkgCards.CreateParams{
				Title:   "Lab 1",
				Content: "Надо сделать",
				ListID:  27,
			},
			card: models.Card{
				ID:       21,
				ListID:   27,
				Title:    "Lab 1",
				Content:  "Надо сделать",
				Position: 41,
			},
			err: nil,
		},
		"list not found": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Create(f.params).Return(*f.card, pkgErrors.ErrListNotFound)
			},
			params: &pkgCards.CreateParams{Title: "Lab 1", Content: "Надо сделать", ListID: 27},
			card:   models.Card{},
			err:    pkgErrors.ErrListNotFound,
		},
		"storages error": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Create(f.params).Return(*f.card, pkgErrors.ErrDb)
			},
			params: &pkgCards.CreateParams{Title: "Lab 1", Content: "Надо сделать", ListID: 27},
			card:   models.Card{},
			err:    pkgErrors.ErrDb,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl), params: test.params, card: &test.card}
			if test.prepare != nil {
				test.prepare(&f)
			}

			uc := NewUsecase(f.repo)
			card, err := uc.Create(test.params)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if card != test.card {
				t.Errorf("\nExpected: %v\nGot: %v", test.card, card)
			}
		})
	}
}

func TestUsecase_List(t *testing.T) {
	type fields struct {
		repo   *mocks.MockRepository
		listID int
		cards  []models.Card
	}

	type testCase struct {
		prepare func(f *fields)
		listID  int
		cards   []models.Card
		err     error
	}

	tests := map[string]testCase{
		"normal": {
			prepare: func(f *fields) {
				f.repo.EXPECT().List(f.listID).Return(f.cards, nil)
			},
			listID: 27,
			cards: []models.Card{
				{ID: 21, ListID: 27, Title: "Lab 1", Content: "Надо сделать", Position: 41},
				{ID: 22, ListID: 27, Title: "Lab 2", Content: "Надо сделать", Position: 42},
				{ID: 23, ListID: 27, Title: "Theory", Content: "Надо выучить", Position: 43},
			},
			err: nil,
		},
		"empty result": {
			prepare: func(f *fields) {
				f.repo.EXPECT().List(f.listID).Return(f.cards, nil)
			},
			listID: 27,
			cards:  []models.Card{},
			err:    nil,
		},
		"list not found": {
			prepare: func(f *fields) {
				f.repo.EXPECT().List(f.listID).Return(f.cards, pkgErrors.ErrListNotFound)
			},
			listID: 27,
			cards:  nil,
			err:    pkgErrors.ErrListNotFound,
		},
		"storages error": {
			prepare: func(f *fields) {
				f.repo.EXPECT().List(f.listID).Return(f.cards, pkgErrors.ErrDb)
			},
			listID: 27,
			cards:  nil,
			err:    pkgErrors.ErrDb,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl), listID: test.listID, cards: test.cards}
			if test.prepare != nil {
				test.prepare(&f)
			}

			serv := NewUsecase(f.repo)
			cards, err := serv.List(test.listID)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if !reflect.DeepEqual(cards, test.cards) {
				t.Errorf("\nExpected: %v\nGot: %v", test.cards, cards)
			}
		})
	}
}

func TestUsecase_Get(t *testing.T) {
	type fields struct {
		repo *mocks.MockRepository
		id   int
		card *models.Card
	}

	type testCase struct {
		prepare func(f *fields)
		id      int
		card    models.Card
		err     error
	}

	tests := map[string]testCase{
		"normal": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Get(f.id).Return(*f.card, nil)
			},
			id:   21,
			card: models.Card{ID: 21, ListID: 27, Title: "Lab 1", Content: "Надо сделать", Position: 41},
			err:  nil,
		},
		"card not found": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Get(f.id).Return(*f.card, pkgErrors.ErrCardNotFound)
			},
			id:   21,
			card: models.Card{},
			err:  pkgErrors.ErrCardNotFound,
		},
		"storages error": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Get(f.id).Return(*f.card, pkgErrors.ErrDb)
			},
			id:   21,
			card: models.Card{},
			err:  pkgErrors.ErrDb,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl), id: test.id, card: &test.card}
			if test.prepare != nil {
				test.prepare(&f)
			}

			uc := NewUsecase(f.repo)
			card, err := uc.Get(test.id)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if card != test.card {
				t.Errorf("\nExpected: %v\nGot: %v", test.card, card)
			}
		})
	}
}

func TestFullUpdate(t *testing.T) {
	type fields struct {
		repo   *mocks.MockRepository
		params *pkgCards.FullUpdateParams
		card   *models.Card
	}

	type testCase struct {
		prepare func(f *fields)
		params  *pkgCards.FullUpdateParams
		card    models.Card
		err     error
	}

	tests := map[string]testCase{
		"normal": {
			prepare: func(f *fields) {
				f.repo.EXPECT().FullUpdate(f.params).Return(*f.card, nil)
			},
			params: &pkgCards.FullUpdateParams{
				ID:       21,
				Title:    "Lab 1",
				Content:  "Надо сделать",
				Position: 41,
				ListID:   27,
			},
			card: models.Card{ID: 21, ListID: 27, Title: "Lab 1", Content: "Надо сделать", Position: 41},
			err:  nil,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl), params: test.params, card: &test.card}
			if test.prepare != nil {
				test.prepare(&f)
			}

			uc := NewUsecase(f.repo)
			card, err := uc.FullUpdate(test.params)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if card != test.card {
				t.Errorf("\nExpected: %v\nGot: %v", test.card, card)
			}
		})
	}
}

func TestPartialUpdate(t *testing.T) {
	type fields struct {
		repo   *mocks.MockRepository
		params *pkgCards.PartialUpdateParams
		card   *models.Card
	}

	type testCase struct {
		prepare func(f *fields)
		params  *pkgCards.PartialUpdateParams
		card    models.Card
		err     error
	}

	tests := map[string]testCase{
		"normal": {
			prepare: func(f *fields) {
				f.repo.EXPECT().PartialUpdate(f.params).Return(*f.card, nil)
			},
			params: &pkgCards.PartialUpdateParams{
				ID:             21,
				Title:          "Lab 1",
				UpdateTitle:    true,
				Content:        "Надо сделать",
				UpdateContent:  true,
				Position:       41,
				UpdatePosition: true,
				ListID:         27,
				UpdateListID:   true,
			},
			card: models.Card{
				ID:       21,
				ListID:   27,
				Title:    "Lab 1",
				Content:  "Надо сделать",
				Position: 41,
			},
			err: nil,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl), params: test.params, card: &test.card}
			if test.prepare != nil {
				test.prepare(&f)
			}

			uc := NewUsecase(f.repo)
			card, err := uc.PartialUpdate(test.params)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if card != test.card {
				t.Errorf("\nExpected: %v\nGot: %v", test.card, card)
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
		"card not found": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Delete(f.id).Return(pkgErrors.ErrCardNotFound)
			},
			id:  21,
			err: pkgErrors.ErrCardNotFound,
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
