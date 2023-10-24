package postgres

import (
	"database/sql"
	pkgCards "github.com/SlavaShagalov/my-trello-backend/internal/cards"
	"github.com/SlavaShagalov/my-trello-backend/internal/models"
	"github.com/SlavaShagalov/my-trello-backend/internal/pkg/constants"
	pkgErrors "github.com/SlavaShagalov/my-trello-backend/internal/pkg/errors"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type repository struct {
	db  *sql.DB
	log *zap.Logger
}

func New(db *sql.DB, log *zap.Logger) pkgCards.Repository {
	return &repository{db: db, log: log}
}

const createCmd = `
	INSERT INTO cards (list_id, title, content, position)
	VALUES ($1, $2, $3, (SELECT COALESCE(MAX(position), 0) + 1
							FROM cards
						  	WHERE list_id = $1))
	RETURNING id, list_id, title, content, position, created_at, updated_at;`

func (repo *repository) Create(params *pkgCards.CreateParams) (models.Card, error) {
	row := repo.db.QueryRow(createCmd, params.ListID, params.Title, params.Content)

	var card models.Card
	err := scanCard(row, &card)
	if err != nil {
		pgErr, ok := err.(*pq.Error)
		if !ok {
			repo.log.Error("Cannot convert err to pq.Error", zap.Error(err))
			return models.Card{}, errors.Wrap(pkgErrors.ErrDb, err.Error())
		}
		if pgErr.Constraint == "cards_list_id_fkey" || pgErr.Code == "23502" {
			return models.Card{}, errors.Wrap(pkgErrors.ErrListNotFound, err.Error())
		}

		repo.log.Error(constants.DBScanError, zap.Error(err), zap.String("sql_query", createCmd),
			zap.Any("create_params", params))
		return models.Card{}, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	repo.log.Debug("New card created", zap.Any("card", card))
	return card, nil
}

const listCmd = `
	SELECT id, list_id, title, content, position, created_at, updated_at
	FROM cards
	WHERE list_id = $1
	ORDER BY position;`

func (repo *repository) ListByList(listID int) ([]models.Card, error) {
	rows, err := repo.db.Query(listCmd, listID)
	if err != nil {
		repo.log.Error(constants.DBError, zap.Error(err), zap.String("sql_query", listCmd),
			zap.Int("list_id", listID))
		return nil, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}
	defer func() {
		_ = rows.Close()
	}()

	cards := []models.Card{}
	var card models.Card
	var content sql.NullString
	for rows.Next() {
		err = rows.Scan(
			&card.ID,
			&card.ListID,
			&card.Title,
			&content,
			&card.Position,
			&card.CreatedAt,
			&card.UpdatedAt,
		)
		if err != nil {
			repo.log.Error(constants.DBScanError, zap.Error(err), zap.String("sql_query", listCmd),
				zap.Int("list_id", listID))
			return nil, errors.Wrap(pkgErrors.ErrDb, err.Error())
		}

		card.Content = content.String
		cards = append(cards, card)
	}

	return cards, nil
}

const listByTitleCmd = `
	SELECT c.id, c.list_id, c.title, c.content, c.position, c.created_at, c.updated_at
	FROM cards c
	JOIN lists l on l.id = c.list_id
	JOIN boards b on b.id = l.board_id
	JOIN workspaces w on w.id = b.workspace_id
	WHERE lower(c.title) LIKE lower('%' || $1 || '%') AND w.user_id = $2
	ORDER BY position;`

func (repo *repository) ListByTitle(title string, listID int) ([]models.Card, error) {
	rows, err := repo.db.Query(listByTitleCmd, title, listID)
	if err != nil {
		repo.log.Error(constants.DBError, zap.Error(err), zap.String("sql", listByTitleCmd),
			zap.String("title", title), zap.Int("list_id", listID))
		return nil, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}
	defer func() {
		_ = rows.Close()
	}()

	cards := []models.Card{}
	var card models.Card
	var content sql.NullString
	for rows.Next() {
		err = rows.Scan(
			&card.ID,
			&card.ListID,
			&card.Title,
			&content,
			&card.Position,
			&card.CreatedAt,
			&card.UpdatedAt,
		)
		if err != nil {
			repo.log.Error(constants.DBScanError, zap.Error(err), zap.String("sql", listByTitleCmd),
				zap.Int("list_id", listID))
			return nil, errors.Wrap(pkgErrors.ErrDb, err.Error())
		}

		card.Content = content.String
		cards = append(cards, card)
	}

	return cards, nil
}

const getCmd = `
	SELECT id, list_id, title, content, position, created_at, updated_at
	FROM cards
	WHERE id = $1;`

func (repo *repository) Get(id int) (models.Card, error) {
	row := repo.db.QueryRow(getCmd, id)

	var card models.Card
	err := scanCard(row, &card)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Card{}, errors.Wrap(pkgErrors.ErrCardNotFound, err.Error())
		}

		repo.log.Error(constants.DBScanError, zap.Error(err), zap.String("sql_query", getCmd),
			zap.Int("id", id))
		return models.Card{}, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	return card, nil
}

const fullUpdateCmd = `
	UPDATE cards
	SET title    = $1,
	    content  = $2,
		position = $3,
		list_id  = $4
	WHERE id = $5
	RETURNING id, list_id, title, content, position, created_at, updated_at;`

func (repo *repository) FullUpdate(params *pkgCards.FullUpdateParams) (models.Card, error) {
	row := repo.db.QueryRow(fullUpdateCmd, params.Title, params.Content, params.Position, params.ListID, params.ID)

	var card models.Card
	err := scanCard(row, &card)
	if err != nil {
		repo.log.Error(constants.DBScanError, zap.Error(err), zap.String("sql_query", fullUpdateCmd),
			zap.Any("params", params))
		return models.Card{}, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	repo.log.Debug("Card full updated", zap.Any("card", card))
	return card, nil
}

const partialUpdateCmd = `
	UPDATE cards
	SET title    = CASE WHEN $1::boolean THEN $2 ELSE title END,
		content  = CASE WHEN $3::boolean THEN $4 ELSE content END,
		position = CASE WHEN $5::boolean THEN $6 ELSE position END,
		list_id  = CASE WHEN $7::boolean THEN $8 ELSE list_id END
	WHERE id = $9
	RETURNING id, list_id, title, content, position, created_at, updated_at;`

const partialUpdateAfterCmd = `
	CALL update_cards_positions($1, $2);`

func (repo *repository) PartialUpdate(params *pkgCards.PartialUpdateParams) (models.Card, error) {
	if params.UpdatePosition {
		_, err := repo.db.Exec(partialUpdateAfterCmd, params.Position, params.ID)
		if err != nil {
			return models.Card{}, err
		}
	}

	row := repo.db.QueryRow(partialUpdateCmd,
		params.UpdateTitle,
		params.Title,
		params.UpdateContent,
		params.Content,
		params.UpdatePosition,
		params.Position,
		params.UpdateListID,
		params.ListID,
		params.ID,
	)

	var card models.Card
	err := scanCard(row, &card)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Card{}, pkgErrors.ErrCardNotFound
		}

		pgErr, _ := err.(*pq.Error)
		if pgErr.Constraint == "cards_list_id_fkey" || pgErr.Code == "23502" {
			return models.Card{}, errors.Wrap(pkgErrors.ErrListNotFound, err.Error())
		}

		repo.log.Error(constants.DBScanError, zap.Error(err), zap.String("sql_query", partialUpdateCmd),
			zap.Any("params", params))
		return models.Card{}, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	repo.log.Debug("Card partial updated", zap.Any("card", card))
	return card, nil
}

const deleteCmd = `
	DELETE FROM cards 
	WHERE id = $1;`

func (repo *repository) Delete(id int) error {
	result, err := repo.db.Exec(deleteCmd, id)
	if err != nil {
		repo.log.Error(constants.DBError, zap.Error(err), zap.String("sql_query", deleteCmd),
			zap.Int("id", id))
		return errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		repo.log.Error(constants.DBError, zap.Error(err), zap.String("sql_query", deleteCmd),
			zap.Int("id", id))
		return errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	if rowsAffected == 0 {
		return pkgErrors.ErrCardNotFound
	}

	repo.log.Debug("Card deleted", zap.Int("id", id))
	return nil
}

func scanCard(row *sql.Row, card *models.Card) error {
	var content sql.NullString
	err := row.Scan(
		&card.ID,
		&card.ListID,
		&card.Title,
		&content,
		&card.Position,
		&card.CreatedAt,
		&card.UpdatedAt,
	)
	if err != nil {
		return err
	}

	card.Content = content.String
	return nil
}
