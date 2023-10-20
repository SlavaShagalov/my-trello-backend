package http

import (
	mw "github.com/SlavaShagalov/my-trello-backend/internal/middleware"
	"github.com/SlavaShagalov/my-trello-backend/internal/pkg/constants"
	pErrors "github.com/SlavaShagalov/my-trello-backend/internal/pkg/errors"
	pHTTP "github.com/SlavaShagalov/my-trello-backend/internal/pkg/http"
	pWorkspaces "github.com/SlavaShagalov/my-trello-backend/internal/workspaces"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type delivery struct {
	uc  pWorkspaces.Usecase
	log *zap.Logger
}

func RegisterHandlers(mux *mux.Router, uc pWorkspaces.Usecase, log *zap.Logger, checkAuth mw.Middleware) {
	del := delivery{
		uc:  uc,
		log: log,
	}

	const (
		workspacesPrefix = "/workspaces"
		workspacesPath   = constants.ApiPrefix + workspacesPrefix
		workspacePath    = workspacesPath + "/{id}"
	)

	mux.HandleFunc(workspacesPath, checkAuth(del.create)).Methods(http.MethodPost)
	mux.HandleFunc(workspacesPath, checkAuth(del.list)).Methods(http.MethodGet)

	mux.HandleFunc(workspacePath, checkAuth(del.get)).Methods(http.MethodGet)
	mux.HandleFunc(workspacePath, checkAuth(del.partialUpdate)).Methods(http.MethodPatch)
	mux.HandleFunc(workspacePath, checkAuth(del.delete)).Methods(http.MethodDelete)
}

func (del *delivery) create(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(mw.ContextUserID).(int)
	if !ok {
		pHTTP.HandleError(w, r, pErrors.ErrReadBody)
		return
	}

	body, err := pHTTP.ReadBody(r, del.log)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	var request createRequest
	err = request.UnmarshalJSON(body)
	if err != nil {
		pHTTP.HandleError(w, r, pErrors.ErrReadBody)
		return
	}

	params := pWorkspaces.CreateParams{
		Title:       request.Title,
		Description: request.Description,
		UserID:      userID,
	}

	workspace, err := del.uc.Create(&params)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	response := newCreateResponse(&workspace)
	pHTTP.SendJSON(w, r, http.StatusOK, response)
}

func (del *delivery) list(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(mw.ContextUserID).(int)
	if !ok {
		pHTTP.HandleError(w, r, pErrors.ErrReadBody)
		return
	}

	workspaces, err := del.uc.List(userID)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	response := newListResponse(workspaces)
	pHTTP.SendJSON(w, r, http.StatusOK, response)
}

func (del *delivery) get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	workspaceID, err := strconv.Atoi(vars["id"])
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	workspace, err := del.uc.Get(workspaceID)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	response := newGetResponse(&workspace)
	pHTTP.SendJSON(w, r, http.StatusOK, response)
}

func (del *delivery) partialUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	workspaceID, err := strconv.Atoi(vars["id"])
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	body, err := pHTTP.ReadBody(r, del.log)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	var request partialUpdateRequest
	err = request.UnmarshalJSON(body)
	if err != nil {
		pHTTP.HandleError(w, r, pErrors.ErrReadBody)
		return
	}

	params := pWorkspaces.PartialUpdateParams{ID: workspaceID}
	params.UpdateTitle = request.Title != nil
	if params.UpdateTitle {
		params.Title = *request.Title
	}
	params.UpdateDescription = request.Description != nil
	if params.UpdateDescription {
		params.Description = *request.Description
	}

	workspace, err := del.uc.PartialUpdate(&params)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	response := newGetResponse(&workspace)
	pHTTP.SendJSON(w, r, http.StatusOK, response)
}

func (del *delivery) delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	workspaceID, err := strconv.Atoi(vars["id"])
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	err = del.uc.Delete(workspaceID)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
