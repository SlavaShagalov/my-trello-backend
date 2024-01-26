package http

import (
	"context"
	"github.com/SlavaShagalov/my-trello-backend/internal/auth"
	mw "github.com/SlavaShagalov/my-trello-backend/internal/middleware"
	"github.com/SlavaShagalov/my-trello-backend/internal/pkg/constants"
	pErrors "github.com/SlavaShagalov/my-trello-backend/internal/pkg/errors"
	pHTTP "github.com/SlavaShagalov/my-trello-backend/internal/pkg/http"
	"github.com/SlavaShagalov/my-trello-backend/internal/pkg/opentel"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
	"net/http"
	"time"
)

const (
	authPrefix = "/auth"
	signInPath = constants.ApiPrefix + authPrefix + "/signin"
	signUpPath = constants.ApiPrefix + authPrefix + "/signup"
	logoutPath = constants.ApiPrefix + authPrefix + "/logout"
)

type delivery struct {
	uc  auth.Usecase
	log *zap.Logger
}

func RegisterHandlers(mux *mux.Router, uc auth.Usecase, log *zap.Logger, checkAuth mw.Middleware, metrics mw.Middleware) {
	del := delivery{
		uc:  uc,
		log: log,
	}

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
	ctx, span := opentel.Tracer.Start(r.Context(), r.Method+" "+signUpPath)
	defer span.End()

	body, err := pHTTP.ReadBody(r, d.log)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	var request SignUpRequest
	err = request.UnmarshalJSON(body)
	if err != nil {
		pHTTP.HandleError(w, r, pErrors.ErrReadBody)
		span.SetStatus(codes.Error, "UnmarshalJSON failed")
		span.RecordError(err)
		return
	}

	params := auth.SignUpParams{
		Name:     request.Name,
		Username: request.Username,
		Email:    request.Email,
		Password: request.Password,
	}

	user, authToken, err := d.uc.SignUp(ctx, &params)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		span.SetStatus(codes.Error, "SignUp failed")
		span.RecordError(err)
		return
	}

	sessionCookie := createSessionCookie(authToken)
	http.SetCookie(w, sessionCookie)
	span.AddEvent("Session cookie created")

	response := newSignUpResponse(&user)
	span.AddEvent("SignUp response created")
	pHTTP.SendJSON(w, r, http.StatusOK, response)
	span.SetStatus(codes.Ok, "SignUp successfully")
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
	ctx, span := opentel.Tracer.Start(r.Context(), "Measure Span")
	_, spanN := opentel.Tracer.Start(context.Background(), r.Method+" "+signInPath)
	defer func() {
		spanN.End()
		span.End()
	}()

	opentel.Counter.Add(ctx, 1)

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
	ctx, span := opentel.Tracer.Start(r.Context(), r.Method+" "+logoutPath)
	defer span.End()

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

	err := d.uc.Logout(ctx, userID, authToken)
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
