package errors

import "net/http"

var httpCodes = map[error]int{
	// Common repository
	ErrDb: http.StatusInternalServerError,

	// Users
	ErrUserNotFound:      http.StatusNotFound,
	ErrUserAlreadyExists: http.StatusConflict,
	ErrTooShortUsername:  http.StatusBadRequest,
	ErrTooLongUsername:   http.StatusBadRequest,
	ErrEmptyName:         http.StatusBadRequest,
	ErrTooLongName:       http.StatusBadRequest,

	// Workspaces
	ErrWorkspaceNotFound: http.StatusNotFound,

	// Boards
	ErrBoardNotFound: http.StatusNotFound,

	// Lists
	ErrListNotFound:           http.StatusNotFound,
	ErrTooLongListTitle:       http.StatusBadRequest,
	ErrTooLongListDescription: http.StatusBadRequest,

	// Cards
	ErrCardNotFound: http.StatusNotFound,

	// Auth
	ErrWrongLoginOrPassword: http.StatusBadRequest,
	ErrSessionNotFound:      http.StatusNotFound,

	// HTTP
	ErrReadBody:         http.StatusBadRequest,
	ErrBadSessionCookie: http.StatusBadRequest,
}

func GetHTTPCodeByError(err error) (int, bool) {
	httpCode, exist := httpCodes[err]
	if !exist {
		httpCode = http.StatusInternalServerError
	}
	return httpCode, exist
}
