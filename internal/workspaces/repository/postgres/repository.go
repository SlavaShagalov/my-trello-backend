package postgres

import (
	"database/sql"
	"github.com/SlavaShagalov/my-trello-backend/internal/models"
	"github.com/SlavaShagalov/my-trello-backend/internal/pkg/constants"
	pkgErrors "github.com/SlavaShagalov/my-trello-backend/internal/pkg/errors"
	pkgWorkspaces "github.com/SlavaShagalov/my-trello-backend/internal/workspaces"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type repository struct {
	db  *sql.DB
	log *zap.Logger
}

func New(db *sql.DB, log *zap.Logger) pkgWorkspaces.Repository {
	return &repository{db: db, log: log}
}

const createCmd = `
	INSERT INTO workspaces (user_id, title, description) 
	VALUES ($1, $2, $3)
	RETURNING *;`

func (repo *repository) Create(params *pkgWorkspaces.CreateParams) (models.Workspace, error) {
	row := repo.db.QueryRow(createCmd, params.UserID, params.Title, params.Description)

	var workspace models.Workspace
	err := scanWorkspace(row, &workspace)
	if err != nil {
		pgErr, ok := err.(*pq.Error)
		if !ok {
			repo.log.Error("Cannot convert err to pq.Error", zap.Error(err))
			return models.Workspace{}, errors.Wrap(pkgErrors.ErrDb, err.Error())
		}
		if pgErr.Constraint == "workspaces_user_id_fkey" {
			return models.Workspace{}, errors.Wrap(pkgErrors.ErrUserNotFound, err.Error())
		}

		repo.log.Error(constants.DBScanError, zap.Error(err), zap.String("sql_query", createCmd),
			zap.Any("create_params", params))
		return models.Workspace{}, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	repo.log.Debug("New workspace", zap.Int("workspace_id", workspace.ID))
	return workspace, nil
}

const listCmd = `
	SELECT id, user_id, title, description, created_at, updated_at
	FROM workspaces
	WHERE user_id = $1;`

func (repo *repository) List(userID int) ([]models.Workspace, error) {
	rows, err := repo.db.Query(listCmd, userID)
	if err != nil {
		repo.log.Error(constants.DBQueryError, zap.Error(err), zap.String("sql_query", listCmd),
			zap.Int("user_id", userID))
		return nil, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}
	defer func() {
		_ = rows.Close()
	}()

	workspaces := []models.Workspace{}
	var workspace models.Workspace
	var description sql.NullString
	for rows.Next() {
		err = rows.Scan(
			&workspace.ID,
			&workspace.UserID,
			&workspace.Title,
			&description,
			&workspace.CreatedAt,
			&workspace.UpdatedAt,
		)
		if err != nil {
			repo.log.Error(constants.DBScanError, zap.Error(err), zap.String("sql_query", listCmd),
				zap.Int("user_id", userID))
			return nil, errors.Wrap(pkgErrors.ErrDb, err.Error())
		}

		workspace.Description = description.String
		workspaces = append(workspaces, workspace)
	}

	return workspaces, nil
}

const getCmd = `
	SELECT id, user_id, title, description, created_at, updated_at
	FROM workspaces
	WHERE id = $1;`

func (repo *repository) Get(id int) (models.Workspace, error) {
	row := repo.db.QueryRow(getCmd, id)

	var workspace models.Workspace
	err := scanWorkspace(row, &workspace)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Workspace{}, errors.Wrap(pkgErrors.ErrWorkspaceNotFound, err.Error())
		}

		repo.log.Error(constants.DBScanError, zap.Error(err), zap.String("sql_query", getCmd),
			zap.Int("id", id))
		return models.Workspace{}, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	return workspace, nil
}

const fullUpdateCmd = `
	UPDATE workspaces
	SET title       = $1,
		description = $2
	WHERE id = $3
	RETURNING id, user_id, title, description, created_at, updated_at;`

func (repo *repository) FullUpdate(params *pkgWorkspaces.FullUpdateParams) (models.Workspace, error) {
	row := repo.db.QueryRow(fullUpdateCmd, params.Title, params.Description, params.ID)

	var workspace models.Workspace
	err := scanWorkspace(row, &workspace)
	if err != nil {
		repo.log.Error(constants.DBScanError, zap.Error(err), zap.String("sql_query", fullUpdateCmd),
			zap.Any("params", params))
		return models.Workspace{}, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	repo.log.Debug("Workspace full updated", zap.Any("workspace", workspace))
	return workspace, nil
}

const partialUpdateCmd = `
	UPDATE workspaces
	SET title       = CASE WHEN $1::boolean THEN $2 ELSE title END,
		description = CASE WHEN $3::boolean THEN $4 ELSE description END
	WHERE id = $5
	RETURNING id, user_id, title, description, created_at, updated_at;`

func (repo *repository) PartialUpdate(params *pkgWorkspaces.PartialUpdateParams) (models.Workspace, error) {
	row := repo.db.QueryRow(partialUpdateCmd,
		params.UpdateTitle,
		params.Title,
		params.UpdateDescription,
		params.Description,
		params.ID,
	)

	var workspace models.Workspace
	err := scanWorkspace(row, &workspace)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Workspace{}, pkgErrors.ErrWorkspaceNotFound
		}

		repo.log.Error(constants.DBScanError, zap.Error(err), zap.String("sql_query", partialUpdateCmd),
			zap.Any("params", params))
		return models.Workspace{}, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	repo.log.Debug("Workspace partial updated", zap.Any("workspace", workspace))
	return workspace, nil
}

const deleteCmd = `
	DELETE FROM workspaces 
	WHERE id = $1;`

func (repo *repository) Delete(id int) error {
	result, err := repo.db.Exec(deleteCmd, id)
	if err != nil {
		repo.log.Error(constants.DBQueryError, zap.Error(err), zap.String("sql_query", deleteCmd),
			zap.Int("id", id))
		return errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		repo.log.Error(constants.DBQueryError, zap.Error(err), zap.String("sql_query", deleteCmd),
			zap.Int("id", id))
		return errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	if rowsAffected == 0 {
		return pkgErrors.ErrWorkspaceNotFound
	}

	repo.log.Debug("Workspace deleted", zap.Int("id", id))
	return nil
}

func scanWorkspace(row *sql.Row, workspace *models.Workspace) error {
	var description sql.NullString
	err := row.Scan(
		&workspace.ID,
		&workspace.UserID,
		&workspace.Title,
		&description,
		&workspace.CreatedAt,
		&workspace.UpdatedAt,
	)
	if err != nil {
		return err
	}

	workspace.Description = description.String
	return nil
}
