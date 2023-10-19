package http

import (
	"github.com/SlavaShagalov/my-trello-backend/internal/auth"
	mw "github.com/SlavaShagalov/my-trello-backend/internal/middleware"
	"github.com/SlavaShagalov/my-trello-backend/internal/pkg/constants"
	pErrors "github.com/SlavaShagalov/my-trello-backend/internal/pkg/errors"
	pHTTP "github.com/SlavaShagalov/my-trello-backend/internal/pkg/http"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type delivery struct {
	uc  auth.Usecase
	log *zap.Logger
}

func RegisterHandlers(mux *mux.Router, uc auth.Usecase, log *zap.Logger, checkAuth mw.Middleware) {
	del := delivery{
		uc:  uc,
		log: log,
	}

	const (
		authPrefix = "/auth"
		signInPath = constants.ApiPrefix + authPrefix + "/signin"
		signUpPath = constants.ApiPrefix + authPrefix + "/signup"
		logoutPath = constants.ApiPrefix + authPrefix + "/logout"
	)

	mux.HandleFunc(signUpPath, del.signup).Methods(http.MethodPost)
	mux.HandleFunc(signInPath, del.signin).Methods(http.MethodPost)
	mux.HandleFunc(logoutPath, checkAuth(del.logout)).Methods(http.MethodDelete)
}

func (del *delivery) signup(w http.ResponseWriter, r *http.Request) {
	body, err := pHTTP.ReadBody(r, del.log)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	var request signUpRequest
	err = request.UnmarshalJSON(body)
	if err != nil {
		pHTTP.HandleError(w, r, pErrors.ErrReadBody)
		return
	}

	params := auth.SignUpParams{
		Name:     request.Name,
		Username: request.Username,
		Email:    request.Email,
		Password: request.Password,
	}

	user, authToken, err := del.uc.SignUp(&params)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	sessionCookie := createSessionCookie(authToken)
	http.SetCookie(w, sessionCookie)

	response := newSignUpResponse(&user)
	pHTTP.SendJSON(w, r, http.StatusOK, response)
}

func (del *delivery) signin(w http.ResponseWriter, r *http.Request) {
	body, err := pHTTP.ReadBody(r, del.log)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	var request signInRequest
	err = request.UnmarshalJSON(body)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	params := auth.SignInParams{
		Username: request.Username,
		Password: request.Password,
	}

	user, authToken, err := del.uc.SignIn(&params)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	sessionCookie := createSessionCookie(authToken)
	http.SetCookie(w, sessionCookie)

	response := newSignInResponse(&user)
	pHTTP.SendJSON(w, r, http.StatusOK, response)
}

func (del *delivery) logout(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(mw.ContextUserID).(int)
	if !ok {
		pHTTP.HandleError(w, r, pErrors.ErrReadBody)
		return
	}
	authToken, ok := r.Context().Value(mw.ContextAuthToken).(string)
	if !ok {
		pHTTP.HandleError(w, r, pErrors.ErrReadBody)
		return
	}

	err := del.uc.Logout(userID, authToken)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	newCookie := &http.Cookie{
		Name:     constants.SessionName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Now().Add(-24 * time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(w, newCookie)

	w.WriteHeader(http.StatusNoContent)
}
