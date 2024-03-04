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
		Path:     "/",
		Expires:  time.Now().Add(constants.SessionLivingTime),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}
}
