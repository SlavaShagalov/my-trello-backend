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

	pkgCards "github.com/SlavaShagalov/my-trello-backend/internal/cards"
	pkgZap "github.com/SlavaShagalov/my-trello-backend/internal/pkg/log/zap"
	pkgDb "github.com/SlavaShagalov/my-trello-backend/internal/pkg/storages"

	cardsRepo "github.com/SlavaShagalov/my-trello-backend/internal/cards/repository/postgres"
	cardsUC "github.com/SlavaShagalov/my-trello-backend/internal/cards/usecase"
)

type CardsSuite struct {
	suite.Suite
	db     *sql.DB
	logger *zap.Logger
	uc     pkgCards.Usecase
}

func (s *CardsSuite) SetupSuite() {
	s.logger = pkgZap.NewTestLogger()

	var err error
	config.SetTestPostgresConfig()
	s.db, err = pkgDb.NewPostgres(s.logger)
	s.Require().NoError(err)

	repo := cardsRepo.New(s.db, s.logger)
	s.uc = cardsUC.New(repo)
}

func (s *CardsSuite) TearDownSuite() {
	err := s.db.Close()
	s.Require().NoError(err)

	_ = s.logger.Sync()
}

func (s *CardsSuite) TestCreate() {
	type testCase struct {
		params *pkgCards.CreateParams
		err    error
	}

	tests := map[string]testCase{
		"normal": {
			params: &pkgCards.CreateParams{
				Title:   "Lab 1",
				Content: "Надо сделать",
				ListID:  3,
			},
			err: nil,
		},
		"list not found": {
			params: &pkgCards.CreateParams{
				Title:   "Lab 1",
				Content: "Надо сделать",
				ListID:  999,
			},
			err: pkgErrors.ErrListNotFound,
		},
	}

	for name, test := range tests {
		s.Run(name, func() {
			card, err := s.uc.Create(test.params)
			assert.ErrorIs(s.T(), err, test.err, "unexpected error")

			if err == nil {
				assert.Equal(s.T(), test.params.ListID, card.ListID, "incorrect ListID")
				assert.Equal(s.T(), test.params.Title, card.Title, "incorrect Title")
				assert.Equal(s.T(), test.params.Content, card.Content, "incorrect Content")

				getCard, err := s.uc.Get(card.ID)
				assert.NoError(s.T(), err, "failed to fetch card from the database")
				assert.Equal(s.T(), card.ID, getCard.ID, "incorrect cardID")
				assert.Equal(s.T(), test.params.ListID, getCard.ListID, "incorrect ListID")
				assert.Equal(s.T(), test.params.Title, getCard.Title, "incorrect Title")
				assert.Equal(s.T(), test.params.Content, getCard.Content, "incorrect Content")

				err = s.uc.Delete(card.ID)
				assert.NoError(s.T(), err, "failed to delete created card")
			}
		})
	}
}

func (s *CardsSuite) TestList() {
	type testCase struct {
		userID int
		cards  []models.Card
		err    error
	}

	tests := map[string]testCase{
		"normal": {
			userID: 2,
			cards: []models.Card{
				{
					ID:       4,
					ListID:   2,
					Title:    "Разработка интерфейса",
					Content:  "Разработать пользовательский интерфейс приложения",
					Position: 1,
				},
				{
					ID:       5,
					ListID:   2,
					Title:    "Написание кода",
					Content:  "Написать программный код для реализации функциональности",
					Position: 2,
				},
				{
					ID:       6,
					ListID:   2,
					Title:    "Тестирование модуля",
					Content:  "Провести тестирование разработанного модуля",
					Position: 3,
				},
			},
			err: nil,
		},
		"empty result": {
			userID: 21,
			cards:  []models.Card{},
			err:    nil,
		},
	}

	for name, test := range tests {
		s.Run(name, func() {
			cards, err := s.uc.ListByList(test.userID)

			assert.ErrorIs(s.T(), err, test.err, "unexpected error")

			if err == nil {
				assert.Equal(s.T(), len(test.cards), len(cards), "incorrect cards length")
				for i := 0; i < len(test.cards); i++ {
					assert.Equal(s.T(), test.cards[i].ID, cards[i].ID, "incorrect cardID")
					assert.Equal(s.T(), test.cards[i].ListID, cards[i].ListID, "incorrect ListID")
					assert.Equal(s.T(), test.cards[i].Title, cards[i].Title, "incorrect Title")
					assert.Equal(s.T(), test.cards[i].Content, cards[i].Content, "incorrect Content")
					assert.Equal(s.T(), test.cards[i].Position, cards[i].Position, "incorrect Position")
				}
			}
		})
	}
}

func (s *CardsSuite) TestGet() {
	type testCase struct {
		cardID int
		card   models.Card
		err    error
	}

	tests := map[string]testCase{
		"normal": {
			cardID: 8,
			card: models.Card{
				ID:       8,
				ListID:   3,
				Title:    "Проведение приемочных тестов",
				Content:  "Провести приемочное тестирование и подтвердить работоспособность",
				Position: 2,
			},
			err: nil,
		},
		"card not found": {
			cardID: 999,
			card:   models.Card{},
			err:    pkgErrors.ErrCardNotFound,
		},
	}

	for name, test := range tests {
		s.Run(name, func() {
			card, err := s.uc.Get(test.cardID)

			assert.ErrorIs(s.T(), err, test.err, "unexpected error")

			if err == nil {
				assert.Equal(s.T(), test.card.ID, card.ID, "incorrect cardID")
				assert.Equal(s.T(), test.card.ListID, card.ListID, "incorrect ListID")
				assert.Equal(s.T(), test.card.Title, card.Title, "incorrect Title")
				assert.Equal(s.T(), test.card.Content, card.Content, "incorrect Content")
				assert.Equal(s.T(), test.card.Position, card.Position, "incorrect Position")
			}
		})
	}
}

func (s *CardsSuite) TestFullUpdate() {
	type testCase struct {
		params *pkgCards.FullUpdateParams
		err    error
	}

	tests := map[string]testCase{
		"normal": {
			params: &pkgCards.FullUpdateParams{
				Title:   "Lab 1",
				Content: "Надо сделать",
				ListID:  3,
			},
			err: nil,
		},
	}

	for name, test := range tests {
		s.Run(name, func() {
			tempCard, err := s.uc.Create(&pkgCards.CreateParams{
				Title:   "Temp Card",
				Content: "Temp Card Content",
				ListID:  2,
			})
			require.NoError(s.T(), err, "failed to create temp card")

			test.params.ID = tempCard.ID
			card, err := s.uc.FullUpdate(test.params)
			assert.ErrorIs(s.T(), err, test.err, "unexpected error")

			if err == nil {
				// check updated card
				assert.Equal(s.T(), test.params.ID, card.ID, "incorrect ID")
				assert.Equal(s.T(), test.params.Title, card.Title, "incorrect Title")
				assert.Equal(s.T(), test.params.Content, card.Content, "incorrect Content")
				assert.Equal(s.T(), test.params.ListID, card.ListID, "incorrect ListID")

				// check card in storages
				getCard, err := s.uc.Get(card.ID)
				assert.NoError(s.T(), err, "failed to fetch card from the database")
				assert.Equal(s.T(), card.ID, getCard.ID, "incorrect cardID")
				assert.Equal(s.T(), card.ListID, getCard.ListID, "incorrect ListID")
				assert.Equal(s.T(), card.Title, getCard.Title, "incorrect Title")
				assert.Equal(s.T(), card.Content, getCard.Content, "incorrect Content")
			}

			err = s.uc.Delete(tempCard.ID)
			require.NoError(s.T(), err, "failed to delete temp card")
		})
	}
}

func (s *CardsSuite) TestPartialUpdate() {
	type testCase struct {
		params *pkgCards.PartialUpdateParams
		card   models.Card
		err    error
	}

	tests := map[string]testCase{
		"full update": {
			params: &pkgCards.PartialUpdateParams{
				Title:         "Lab 1",
				UpdateTitle:   true,
				Content:       "Надо сделать",
				UpdateContent: true,
			},
			card: models.Card{
				Title:   "Lab 1",
				Content: "Надо сделать",
				ListID:  2,
			},
			err: nil,
		},
		"only title update": {
			params: &pkgCards.PartialUpdateParams{
				Title:       "Lab 1",
				UpdateTitle: true,
			},
			card: models.Card{
				Title:   "Lab 1",
				Content: "Temp Card Content",
				ListID:  2,
			},
			err: nil,
		},
		"only content update": {
			params: &pkgCards.PartialUpdateParams{
				Content:       "Надо сделать",
				UpdateContent: true,
			},
			card: models.Card{
				Title:   "Temp Card",
				Content: "Надо сделать",
				ListID:  2,
			},
			err: nil,
		},
	}

	for name, test := range tests {
		s.Run(name, func() {
			tempCard, err := s.uc.Create(&pkgCards.CreateParams{
				Title:   "Temp Card",
				Content: "Temp Card Content",
				ListID:  2,
			})
			require.NoError(s.T(), err, "failed to create temp card")

			test.params.ID = tempCard.ID
			card, err := s.uc.PartialUpdate(test.params)
			assert.ErrorIs(s.T(), err, test.err, "unexpected error")

			if err == nil {
				// check updated card
				assert.Equal(s.T(), test.params.ID, card.ID, "incorrect cardID")
				assert.Equal(s.T(), test.card.Title, card.Title, "incorrect Title")
				assert.Equal(s.T(), test.card.Content, card.Content, "incorrect Content")
				assert.Equal(s.T(), test.card.ListID, card.ListID, "incorrect ListID")

				// check card in storages
				getCard, err := s.uc.Get(card.ID)
				assert.NoError(s.T(), err, "failed to fetch card from the database")
				assert.Equal(s.T(), test.card.Title, getCard.Title, "incorrect Title")
				assert.Equal(s.T(), test.card.Content, getCard.Content, "incorrect Content")
				assert.Equal(s.T(), test.card.ListID, getCard.ListID, "incorrect ListID")
			}

			err = s.uc.Delete(tempCard.ID)
			require.NoError(s.T(), err, "failed to delete temp card")
		})
	}
}

func (s *CardsSuite) TestDelete() {
	type testCase struct {
		setupCard func() (models.Card, error)
		err       error
	}

	tests := map[string]testCase{
		"normal": {
			setupCard: func() (models.Card, error) {
				return s.uc.Create(&pkgCards.CreateParams{
					Title:   "Test Card",
					Content: "Test Card Content",
					ListID:  1,
				})
			},
			err: nil,
		},
		"card not found": {
			setupCard: func() (models.Card, error) {
				return models.Card{ID: 999}, nil
			},
			err: pkgErrors.ErrCardNotFound,
		},
	}

	for name, test := range tests {
		s.Run(name, func() {
			card, err := test.setupCard()
			s.Require().NoError(err)

			err = s.uc.Delete(card.ID)
			assert.ErrorIs(s.T(), err, test.err, "unexpected error")

			if test.err == nil {
				_, err = s.uc.Get(card.ID)
				assert.ErrorIs(s.T(), err, pkgErrors.ErrCardNotFound, "card should be deleted")
			}
		})
	}
}

func TestCardSuite(t *testing.T) {
	suite.Run(t, new(CardsSuite))
}
