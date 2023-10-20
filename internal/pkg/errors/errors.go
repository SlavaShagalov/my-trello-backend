package errors

import (
	"errors"
	"fmt"
	"github.com/SlavaShagalov/my-trello-backend/internal/pkg/constants"
)

var (
	// Common repository
	ErrDb = errors.New("database error")

	// Users
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrTooShortUsername  = errors.New(fmt.Sprintf("username must be at least %d characters",
		constants.MinUsernameLen))
	ErrTooLongUsername = errors.New(fmt.Sprintf("username must be no more than %d characters",
		constants.MaxUsernameLen))
	ErrEmptyName   = errors.New("name must not be empty")
	ErrTooLongName = errors.New(fmt.Sprintf("name must be no more than %d characters", constants.MaxNameLen))

	// Workspaces
	ErrWorkspaceNotFound = errors.New("workspace not found")

	// Boards
	ErrBoardNotFound = errors.New("board not found")

	// Lists
	ErrListNotFound     = errors.New("list not found")
	ErrTooLongListTitle = errors.New(fmt.Sprintf("list title must be no more than %d characters",
		constants.MaxListTitleLen))
	ErrTooLongListDescription = errors.New(fmt.Sprintf("list description must be no more than %d characters",
		constants.MaxListDescriptionLen))

	// Cards
	ErrCardNotFound = errors.New("card not found")

	// Auth
	ErrWrongLoginOrPassword = errors.New("wrong login or password")
	ErrGetHashedPassword    = errors.New("get hashed password error")
	// ErrSessionStorage       = errors.New("session storage error")
	ErrSessionNotFound = errors.New("session not found")

	// HTTP
	ErrReadBody         = errors.New("read request body error")
	ErrBadSessionCookie = errors.New("bad session cookie")
)
