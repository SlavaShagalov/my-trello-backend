package http

import (
	"github.com/SlavaShagalov/my-trello-backend/internal/pkg/constants"
	"net/http"
	"time"
)

func createSessionCookie(token string) *http.Cookie {
	return &http.Cookie{
		Name:     constants.SessionName,
		Value:    token,
		Expires:  time.Now().Add(constants.SessionLivingTime),
		HttpOnly: true,
		Path:     "/",
	}
}
