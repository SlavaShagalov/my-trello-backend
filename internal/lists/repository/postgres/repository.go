package postgres

import (
	"database/sql"
	pkgLists "github.com/SlavaShagalov/my-trello-backend/internal/lists"
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

func NewRepository(db *sql.DB, log *zap.Logger) pkgLists.Repository {
	return &repository{db: db, log: log}
}

const createCmd = `
	INSERT INTO lists (board_id, title, position) 
	VALUES ($1, $2, (SELECT COALESCE(MAX(position), 0) + 1
						FROM lists
						WHERE board_id = $1))
	RETURNING id, board_id, title, position, created_at, updated_at;`

func (repo *repository) Create(params *pkgLists.CreateParams) (models.List, error) {
	row := repo.db.QueryRow(createCmd, params.BoardID, params.Title)

	var list models.List
	err := scanList(row, &list)
	if err != nil {
		pgErr, ok := err.(*pq.Error)
		if !ok {
			repo.log.Error("Cannot convert err to pq.Error", zap.Error(err))
			return models.List{}, errors.Wrap(pkgErrors.ErrDb, err.Error())
		}
		if pgErr.Constraint == "lists_board_id_fkey" || pgErr.Code == "23502" {
			return models.List{}, errors.Wrap(pkgErrors.ErrBoardNotFound, err.Error())
		}

		repo.log.Error(constants.DBScanError, zap.Error(err), zap.String("sql_query", createCmd),
			zap.Any("create_params", params))
		return models.List{}, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	repo.log.Debug("New list created", zap.Any("list", list))
	return list, nil
}

const listCmd = `
	SELECT id, board_id, title, position, created_at, updated_at
	FROM lists
	WHERE board_id = $1
	ORDER BY position;`

func (repo *repository) List(boardID int) ([]models.List, error) {
	rows, err := repo.db.Query(listCmd, boardID)
	if err != nil {
		repo.log.Error(constants.DBError, zap.Error(err), zap.String("sql_query", listCmd),
			zap.Int("board_id", boardID))
		return nil, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}
	defer func() {
		_ = rows.Close()
	}()

	lists := []models.List{}
	var list models.List
	for rows.Next() {
		err = rows.Scan(
			&list.ID,
			&list.BoardID,
			&list.Title,
			&list.Position,
			&list.CreatedAt,
			&list.UpdatedAt,
		)
		if err != nil {
			repo.log.Error(constants.DBScanError, zap.Error(err), zap.String("sql_query", listCmd),
				zap.Int("board_id", boardID))
			return nil, errors.Wrap(pkgErrors.ErrDb, err.Error())
		}

		lists = append(lists, list)
	}

	return lists, nil
}

const getCmd = `
	SELECT id, board_id, title, position, created_at, updated_at
	FROM lists
	WHERE id = $1;`

func (repo *repository) Get(id int) (models.List, error) {
	row := repo.db.QueryRow(getCmd, id)

	var list models.List
	err := scanList(row, &list)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.List{}, errors.Wrap(pkgErrors.ErrListNotFound, err.Error())
		}

		repo.log.Error(constants.DBScanError, zap.Error(err), zap.String("sql_query", getCmd),
			zap.Int("id", id))
		return models.List{}, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	return list, nil
}

const fullUpdateCmd = `
	UPDATE lists
	SET title    = $1,
		position = $2,
		board_id = $3
	WHERE id = $4
	RETURNING id, board_id, title, position, created_at, updated_at;`

func (repo *repository) FullUpdate(params *pkgLists.FullUpdateParams) (models.List, error) {
	row := repo.db.QueryRow(fullUpdateCmd, params.Title, params.Position, params.BoardID, params.ID)

	var list models.List
	err := scanList(row, &list)
	if err != nil {
		repo.log.Error(constants.DBScanError, zap.Error(err), zap.String("sql_query", fullUpdateCmd),
			zap.Any("params", params))
		return models.List{}, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	repo.log.Debug("ListByWorkspace full updated", zap.Any("list", list))
	return list, nil
}

const partialUpdateCmd = `
	UPDATE lists
	SET title    = CASE WHEN $1 THEN $2 ELSE title END,
		position = CASE WHEN $3 THEN $4 ELSE position END,
		board_id = CASE WHEN $5 THEN $6 ELSE board_id END
	WHERE id = $7
	RETURNING id, board_id, title, position, created_at, updated_at;`

const partialUpdateAfterCmd = `
	CALL update_list_positions($1, $2);`

func (repo *repository) PartialUpdate(params *pkgLists.PartialUpdateParams) (models.List, error) {
	if params.UpdatePosition {
		_, err := repo.db.Exec(partialUpdateAfterCmd, params.Position, params.ID)
		if err != nil {
			return models.List{}, err
		}
	}

	row := repo.db.QueryRow(partialUpdateCmd,
		params.UpdateTitle,
		params.Title,
		params.UpdatePosition,
		params.Position,
		params.UpdateBoardID,
		params.BoardID,
		params.ID,
	)

	var list models.List
	err := scanList(row, &list)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.List{}, pkgErrors.ErrListNotFound
		}

		pgErr, _ := err.(*pq.Error)
		if pgErr.Constraint == "lists_board_id_fkey" || pgErr.Code == "23502" {
			return models.List{}, errors.Wrap(pkgErrors.ErrBoardNotFound, err.Error())
		}

		repo.log.Error(constants.DBScanError, zap.Error(err), zap.String("sql_query", partialUpdateCmd),
			zap.Any("params", params))
		return models.List{}, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	repo.log.Debug("ListByWorkspace partial updated", zap.Any("list", list))
	return list, nil

	//return convert.ListByWorkspace{}, nil
}

const deleteCmd = `
	DELETE FROM lists
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
		return pkgErrors.ErrListNotFound
	}

	repo.log.Debug("ListByWorkspace deleted", zap.Int("id", id))
	return nil
}

func scanList(row *sql.Row, list *models.List) error {
	return row.Scan(
		&list.ID,
		&list.BoardID,
		&list.Title,
		&list.Position,
		&list.CreatedAt,
		&list.UpdatedAt,
	)
}
