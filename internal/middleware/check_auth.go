package middleware

import (
	"context"
	"github.com/SlavaShagalov/my-trello-backend/internal/auth"
	"github.com/SlavaShagalov/my-trello-backend/internal/pkg/constants"
	pErrors "github.com/SlavaShagalov/my-trello-backend/internal/pkg/errors"
	pHTTP "github.com/SlavaShagalov/my-trello-backend/internal/pkg/http"
	"github.com/SlavaShagalov/my-trello-backend/internal/pkg/opentel"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
)

func NewCheckAuth(uc auth.Usecase, log *zap.Logger) func(h http.HandlerFunc) http.HandlerFunc {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx, span := opentel.Tracer.Start(r.Context(), "CheckAuth Middleware")
			defer span.End()
			r = r.WithContext(ctx)

			sessionCookie, err := r.Cookie(constants.SessionName)
			if err != nil {
				log.Debug("Failed to get session cookie", zap.Error(err))
				pHTTP.HandleError(w, r, pErrors.ErrSessionNotFound)
				return
			}

			id, authToken, err := parseSessionCookie(sessionCookie)
			if err != nil {
				pHTTP.HandleError(w, r, pErrors.ErrBadSessionCookie)
				return
			}

			userID, err := uc.CheckAuth(r.Context(), id, authToken)
			if err != nil {
				pHTTP.HandleError(w, r, err)
				return
			}

			ctx = context.WithValue(ctx, ContextUserID, userID)
			ctx = context.WithValue(ctx, ContextAuthToken, authToken)

			h(w, r.WithContext(ctx))
		}
	}
}

func parseSessionCookie(c *http.Cookie) (int, string, error) {
	tmp := strings.Split(c.Value, "$")
	if len(tmp) != 2 {
		return 0, "", pErrors.ErrBadSessionCookie
	}

	id, err := strconv.Atoi(tmp[0])
	if err != nil {
		return 0, "", err
	}

	return id, c.Value, nil
}
