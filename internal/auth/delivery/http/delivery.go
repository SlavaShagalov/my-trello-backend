package http

import (
	"github.com/SlavaShagalov/my-trello-backend/internal/auth"
	mw "github.com/SlavaShagalov/my-trello-backend/internal/middleware"
	"github.com/SlavaShagalov/my-trello-backend/internal/pkg/constants"
	pErrors "github.com/SlavaShagalov/my-trello-backend/internal/pkg/errors"
	pHTTP "github.com/SlavaShagalov/my-trello-backend/internal/pkg/http"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type delivery struct {
	uc     auth.Usecase
	log    *zap.Logger
	tracer trace.Tracer
}

func RegisterHandlers(mux *mux.Router, uc auth.Usecase, log *zap.Logger, checkAuth mw.Middleware, metrics mw.Middleware,
	tracer trace.Tracer) {
	del := delivery{
		uc:     uc,
		log:    log,
		tracer: tracer,
	}

	const (
		authPrefix = "/auth"
		signInPath = constants.ApiPrefix + authPrefix + "/signin"
		signUpPath = constants.ApiPrefix + authPrefix + "/signup"
		logoutPath = constants.ApiPrefix + authPrefix + "/logout"
	)

	mux.HandleFunc(signUpPath, metrics(del.signup)).Methods(http.MethodPost)
	mux.HandleFunc(signInPath, metrics(del.signin)).Methods(http.MethodPost)
	mux.HandleFunc(logoutPath, metrics(checkAuth(del.logout))).Methods(http.MethodDelete)
}

// signup godoc
//
//	@Summary		Creates new user and returns authentication cookie.
//	@Description	Creates new user and returns authentication cookie.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			signUpParams	body		SignUpRequest	true	"Sign up params."
//	@Success		200				{object}	SignUpResponse	"Successfully created user."
//	@Failure		400				{object}	http.JSONError
//	@Failure		405
//	@Failure		500
//	@Router			/auth/signup [post]
func (d *delivery) signup(w http.ResponseWriter, r *http.Request) {
	body, err := pHTTP.ReadBody(r, d.log)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	var request SignUpRequest
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

	user, authToken, err := d.uc.SignUp(&params)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	sessionCookie := createSessionCookie(authToken)
	http.SetCookie(w, sessionCookie)

	response := newSignUpResponse(&user)
	pHTTP.SendJSON(w, r, http.StatusOK, response)
}

// signin godoc
//
//	@Summary		Logs in and returns the authentication cookie
//	@Description	Logs in and returns the authentication cookie
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			signInParams	body		SignInRequest	true	"Successfully authenticated."
//	@Success		200				{object}	SignInResponse	"successfully auth"
//	@Failure		400				{object}	http.JSONError
//	@Failure		404				{object}	http.JSONError
//	@Failure		405
//	@Failure		500
//	@Router			/auth/signin [post]
func (d *delivery) signin(w http.ResponseWriter, r *http.Request) {
	ctx, span := d.tracer.Start(r.Context(), "HTTP GET /signin")
	time.Sleep(3 * time.Millisecond)
	defer span.End()

	body, err := pHTTP.ReadBody(r, d.log)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	var request SignInRequest
	err = request.UnmarshalJSON(body)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	params := auth.SignInParams{
		Username: request.Username,
		Password: request.Password,
	}

	user, authToken, err := d.uc.SignIn(ctx, &params)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}
	time.Sleep(2 * time.Millisecond)

	sessionCookie := createSessionCookie(authToken)
	http.SetCookie(w, sessionCookie)

	response := newSignInResponse(&user)
	pHTTP.SendJSON(w, r, http.StatusOK, response)
}

// logout godoc
//
//	@Summary		Logs out and deletes the authentication cookie.
//	@Description	Logs out and deletes the authentication cookie.
//	@Tags			auth
//	@Produce		json
//	@Success		204	"Successfully logged out."
//	@Failure		400	{object}	http.JSONError
//	@Failure		401	{object}	http.JSONError
//	@Failure		405
//	@Failure		500
//	@Router			/auth/logout [delete]
//
//	@Security		cookieAuth
func (d *delivery) logout(w http.ResponseWriter, r *http.Request) {
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

	err := d.uc.Logout(userID, authToken)
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
