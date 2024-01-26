package std

import (
	"context"
	"database/sql"
	pkgBoards "github.com/SlavaShagalov/my-trello-backend/internal/boards"
	"github.com/SlavaShagalov/my-trello-backend/internal/models"
	"github.com/SlavaShagalov/my-trello-backend/internal/pkg/constants"
	pkgErrors "github.com/SlavaShagalov/my-trello-backend/internal/pkg/errors"
	"github.com/SlavaShagalov/my-trello-backend/internal/pkg/opentel"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	componentName = "Boards Repository"
)

type repository struct {
	db  *sql.DB
	log *zap.Logger
}

func New(db *sql.DB, log *zap.Logger) pkgBoards.Repository {
	return &repository{db: db, log: log}
}

const createCmd = `
	INSERT INTO boards (workspace_id, title, description) 
	VALUES ($1, $2, $3)
	RETURNING id, workspace_id, title, description, background, created_at, updated_at;`

func (repo *repository) Create(ctx context.Context, params *pkgBoards.CreateParams) (models.Board, error) {
	_, span := opentel.Tracer.Start(ctx, componentName+" "+"Create")
	defer span.End()

	row := repo.db.QueryRow(createCmd, params.WorkspaceID, params.Title, params.Description)

	var board models.Board
	err := scanBoard(row, &board)
	if err != nil {
		pgErr, ok := err.(*pq.Error)
		if !ok {
			repo.log.Error("Cannot convert err to pq.Error", zap.Error(err))
			return models.Board{}, errors.Wrap(pkgErrors.ErrDb, err.Error())
		}
		if pgErr.Constraint == "boards_workspace_id_fkey" {
			return models.Board{}, errors.Wrap(pkgErrors.ErrWorkspaceNotFound, err.Error())
		}

		repo.log.Error(constants.DBScanError, zap.Error(err), zap.String("sql_query", createCmd),
			zap.Any("create_params", params))
		return models.Board{}, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	repo.log.Debug("New board created", zap.Any("board", board))
	return board, nil
}

const listCmd = `
	SELECT id, workspace_id, title, description, background, created_at, updated_at
	FROM boards
	WHERE workspace_id = $1;`

func (repo *repository) List(ctx context.Context, workspaceID int) ([]models.Board, error) {
	_, span := opentel.Tracer.Start(ctx, componentName+" "+"List")
	defer span.End()

	rows, err := repo.db.Query(listCmd, workspaceID)
	if err != nil {
		repo.log.Error(constants.DBError, zap.Error(err), zap.String("sql_query", listCmd),
			zap.Int("workspace_id", workspaceID))
		return nil, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}
	defer func() {
		_ = rows.Close()
	}()

	boards := []models.Board{}
	var board models.Board
	var description sql.NullString
	background := new(sql.NullString)
	for rows.Next() {
		err = rows.Scan(
			&board.ID,
			&board.WorkspaceID,
			&board.Title,
			&description,
			background,
			&board.CreatedAt,
			&board.UpdatedAt,
		)
		if err != nil {
			repo.log.Error(constants.DBScanError, zap.Error(err), zap.String("sql_query", listCmd),
				zap.Int("workspace_id", workspaceID))
			return nil, errors.Wrap(pkgErrors.ErrDb, err.Error())
		}

		if background.Valid {
			board.Background = &background.String
		} else {
			board.Background = nil
		}
		board.Description = description.String

		boards = append(boards, board)
	}

	return boards, nil
}

const listByTitleCmd = `
	SELECT b.id, b.workspace_id, b.title, b.description, b.background, b.created_at, b.updated_at
	FROM boards b 
	JOIN workspaces w on w.id = b.workspace_id
	WHERE lower(b.title) LIKE lower('%' || $1 || '%') AND w.user_id = $2;`

func (repo *repository) ListByTitle(ctx context.Context, title string, userID int) ([]models.Board, error) {
	_, span := opentel.Tracer.Start(ctx, componentName+" "+"ListByTitle")
	defer span.End()

	rows, err := repo.db.Query(listByTitleCmd, title, userID)
	if err != nil {
		repo.log.Error(constants.DBError, zap.Error(err), zap.String("sql", listByTitleCmd),
			zap.String("title", title))
		return nil, pkgErrors.ErrDb
	}
	defer func() {
		_ = rows.Close()
	}()

	boards := []models.Board{}
	var board models.Board
	var description sql.NullString
	background := new(sql.NullString)
	for rows.Next() {
		err = rows.Scan(
			&board.ID,
			&board.WorkspaceID,
			&board.Title,
			&description,
			background,
			&board.CreatedAt,
			&board.UpdatedAt,
		)
		if err != nil {
			repo.log.Error(constants.DBScanError, zap.Error(err), zap.String("sql", listByTitleCmd),
				zap.String("title", title))
			return nil, errors.Wrap(pkgErrors.ErrDb, err.Error())
		}

		if background.Valid {
			board.Background = &background.String
		} else {
			board.Background = nil
		}
		board.Description = description.String

		boards = append(boards, board)
	}

	return boards, nil
}

const getCmd = `
	SELECT id, workspace_id, title, description, background, created_at, updated_at
	FROM boards
	WHERE id = $1;`

func (repo *repository) Get(ctx context.Context, id int) (models.Board, error) {
	_, span := opentel.Tracer.Start(ctx, componentName+" "+"Get")
	defer span.End()

	row := repo.db.QueryRow(getCmd, id)

	var board models.Board
	err := scanBoard(row, &board)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Board{}, errors.Wrap(pkgErrors.ErrBoardNotFound, err.Error())
		}

		repo.log.Error(constants.DBScanError, zap.Error(err), zap.String("sql_query", getCmd),
			zap.Int("id", id))
		return models.Board{}, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	return board, nil
}

const fullUpdateCmd = `
	UPDATE boards
	SET title        = $1,
		description  = $2,
		workspace_id = $3
	WHERE id = $4
	RETURNING id, workspace_id, title, description, background, created_at, updated_at;`

func (repo *repository) FullUpdate(ctx context.Context, params *pkgBoards.FullUpdateParams) (models.Board, error) {
	_, span := opentel.Tracer.Start(ctx, componentName+" "+"FullUpdate")
	defer span.End()

	row := repo.db.QueryRow(fullUpdateCmd, params.Title, params.Description, params.WorkspaceID, params.ID)

	var board models.Board
	err := scanBoard(row, &board)
	if err != nil {
		repo.log.Error(constants.DBScanError, zap.Error(err), zap.String("sql_query", fullUpdateCmd),
			zap.Any("params", params))
		return models.Board{}, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	repo.log.Debug("Board full updated", zap.Any("board", board))
	return board, nil
}

const partialUpdateCmd = `
	UPDATE boards
	SET title        = CASE WHEN $1::boolean THEN $2 ELSE title END,
		description  = CASE WHEN $3::boolean THEN $4 ELSE description END,
		workspace_id = CASE WHEN $5::boolean THEN $6 ELSE workspace_id END
	WHERE id = $7
	RETURNING id, workspace_id, title, description, background, created_at, updated_at;`

func (repo *repository) PartialUpdate(ctx context.Context, params *pkgBoards.PartialUpdateParams) (models.Board, error) {
	_, span := opentel.Tracer.Start(ctx, componentName+" "+"PartialUpdate")
	defer span.End()

	row := repo.db.QueryRow(partialUpdateCmd,
		params.UpdateTitle,
		params.Title,
		params.UpdateDescription,
		params.Description,
		params.UpdateWorkspaceID,
		params.WorkspaceID,
		params.ID,
	)

	var board models.Board
	err := scanBoard(row, &board)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Board{}, errors.Wrap(pkgErrors.ErrBoardNotFound, err.Error())
		}

		repo.log.Error(constants.DBScanError, zap.Error(err), zap.String("sql_query", partialUpdateCmd),
			zap.Any("params", params))
		return models.Board{}, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	repo.log.Debug("Board partial updated", zap.Any("board", board))
	return board, nil
}

const updateBackgroundCmd = `
	UPDATE boards
	SET background = $1
	WHERE id = $2;`

func (repo *repository) UpdateBackground(ctx context.Context, id int, background string) error {
	_, span := opentel.Tracer.Start(ctx, componentName+" "+"UpdateBackground")
	defer span.End()

	result, err := repo.db.Exec(updateBackgroundCmd, background, id)
	if err != nil {
		repo.log.Error(constants.DBError, zap.Error(err), zap.String("sql", updateBackgroundCmd),
			zap.Int("id", id))
		return pkgErrors.ErrDb
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		repo.log.Error(constants.DBError, zap.Error(err), zap.String("sql", updateBackgroundCmd),
			zap.Int("id", id))
		return pkgErrors.ErrDb
	}

	if rowsAffected == 0 {
		return pkgErrors.ErrBoardNotFound
	}

	repo.log.Debug("Background updated", zap.Int("id", id))
	return nil
}

const deleteCmd = `
	DELETE FROM boards 
	WHERE id = $1;`

func (repo *repository) Delete(ctx context.Context, id int) error {
	_, span := opentel.Tracer.Start(ctx, componentName+" "+"Delete")
	defer span.End()

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
		return pkgErrors.ErrBoardNotFound
	}

	repo.log.Debug("Board deleted", zap.Int("id", id))
	return nil
}

func scanBoard(row *sql.Row, board *models.Board) error {
	var description sql.NullString
	background := new(sql.NullString)
	err := row.Scan(
		&board.ID,
		&board.WorkspaceID,
		&board.Title,
		&description,
		background,
		&board.CreatedAt,
		&board.UpdatedAt,
	)
	if err != nil {
		return err
	}

	if background.Valid {
		board.Background = &background.String
	} else {
		board.Background = nil
	}

	board.Description = description.String
	return nil
}
