package postgres

import (
	"database/sql"
	"github.com/SlavaShagalov/my-trello-backend/internal/models"
	"github.com/SlavaShagalov/my-trello-backend/internal/pkg/constants"
	pkgErrors "github.com/SlavaShagalov/my-trello-backend/internal/pkg/errors"
	pkgUsers "github.com/SlavaShagalov/my-trello-backend/internal/users"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type repository struct {
	db  *sql.DB
	log *zap.Logger
}

func New(db *sql.DB, log *zap.Logger) pkgUsers.Repository {
	return &repository{db: db, log: log}
}

const createCmd = `
	INSERT INTO users (name, username, email, hashed_password)
	VALUES ($1, $2, $3, $4)
	RETURNING id, username, hashed_password, email, name, avatar, created_at, updated_at;`

func (repo *repository) Create(params *pkgUsers.CreateParams) (models.User, error) {
	row := repo.db.QueryRow(createCmd, params.Name, params.Username, params.Email, params.HashedPassword)

	var user models.User
	err := scanUser(row, &user)
	if err != nil {
		repo.log.Error(constants.DBScanError, zap.Error(err), zap.String("sql_query", createCmd),
			zap.Any("create_params", params))
		return models.User{}, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	return user, nil
}

const listCmd = `
	SELECT id, username, hashed_password, email, name, avatar, created_at, updated_at
	FROM users;`

func (repo *repository) List() ([]models.User, error) {
	rows, err := repo.db.Query(listCmd)
	if err != nil {
		repo.log.Error(constants.DBQueryError, zap.Error(err), zap.String("sql_query", listCmd))
		return nil, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}
	defer func() {
		_ = rows.Close()
	}()

	users := []models.User{}
	var user models.User
	var avatar sql.NullString
	for rows.Next() {
		err = rows.Scan(
			&user.ID,
			&user.Username,
			&user.Password,
			&user.Email,
			&user.Name,
			&avatar,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			repo.log.Error(constants.DBScanError, zap.Error(err), zap.String("sql_query", listCmd))
			return nil, errors.Wrap(pkgErrors.ErrDb, err.Error())
		}

		user.Avatar = avatar.String
		users = append(users, user)
	}

	return users, nil
}

const getCmd = `
	SELECT id, username, hashed_password, email, name, avatar, created_at, updated_at
	FROM users
	WHERE id = $1;`

func (repo *repository) Get(id int) (models.User, error) {
	row := repo.db.QueryRow(getCmd, id)

	var user models.User
	err := scanUser(row, &user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, errors.Wrap(pkgErrors.ErrUserNotFound, err.Error())
		}

		repo.log.Error(constants.DBScanError, zap.Error(err), zap.String("sql_query", getCmd),
			zap.Int("id", id))
		return models.User{}, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	return user, nil
}

const getByUsernameCmd = `
	SELECT id, username, hashed_password, email, name, avatar, created_at, updated_at
	FROM users
	WHERE username = $1;`

func (repo *repository) GetByUsername(username string) (models.User, error) {
	row := repo.db.QueryRow(getByUsernameCmd, username)

	var user models.User
	err := scanUser(row, &user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, pkgErrors.ErrUserNotFound
		}

		repo.log.Error(constants.DBScanError, zap.Error(err), zap.String("sql_query", getByUsernameCmd),
			zap.String("username", username))
		return models.User{}, pkgErrors.ErrDb
	}

	return user, nil
}

const fullUpdateCmd = `
	UPDATE users
	SET username = $1,
	    email    = $2,
		name     = $3
	WHERE id = $4
	RETURNING id, username, hashed_password, email, name, avatar, created_at, updated_at;`

func (repo *repository) FullUpdate(params *pkgUsers.FullUpdateParams) (models.User, error) {
	row := repo.db.QueryRow(fullUpdateCmd, params.Username, params.Email, params.Name, params.ID)

	var user models.User
	err := scanUser(row, &user)
	if err != nil {
		repo.log.Error(constants.DBScanError, zap.Error(err), zap.String("sql_query", fullUpdateCmd),
			zap.Any("params", params))
		return models.User{}, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	repo.log.Debug("User full updated", zap.Any("user", user))
	return user, nil
}

const partialUpdateCmd = `
	UPDATE users
	SET username = CASE WHEN $1::boolean THEN $2 ELSE username END,
		email    = CASE WHEN $3::boolean THEN $4 ELSE email END,
		name     = CASE WHEN $5::boolean THEN $6 ELSE name END
	WHERE id = $7
	RETURNING id, username, hashed_password, email, name, avatar, created_at, updated_at;`

func (repo *repository) PartialUpdate(params *pkgUsers.PartialUpdateParams) (models.User, error) {
	row := repo.db.QueryRow(partialUpdateCmd,
		params.UpdateUsername,
		params.Username,
		params.UpdateEmail,
		params.Email,
		params.UpdateName,
		params.Name,
		params.ID,
	)

	var user models.User
	err := scanUser(row, &user)
	if err != nil {
		repo.log.Error(constants.DBScanError, zap.Error(err), zap.String("sql_query", partialUpdateCmd),
			zap.Any("params", params))
		return models.User{}, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	repo.log.Debug("User partial updated", zap.Any("user", user))
	return user, nil
}

const deleteCmd = `
	DELETE FROM users 
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
		return pkgErrors.ErrUserNotFound
	}

	repo.log.Debug("User deleted", zap.Int("id", id))
	return nil
}

const existsCmd = `
	SELECT EXISTS(SELECT id
					FROM users
					WHERE id = $1) AS exists;`

func (repo *repository) Exists(userID int) (bool, error) {
	row := repo.db.QueryRow(existsCmd, userID)

	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		repo.log.Error(constants.DBScanError, zap.Error(err), zap.String("sql_query", existsCmd),
			zap.Int("user_id", userID))
		return false, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}
	return exists, nil
}

func scanUser(row *sql.Row, user *models.User) error {
	var avatar sql.NullString
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.Name,
		&avatar,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return err
	}

	user.Avatar = avatar.String
	return nil
}
