package http

import (
	pBoards "github.com/SlavaShagalov/my-trello-backend/internal/boards"
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
	uc       pWorkspaces.Usecase
	boardsUC pBoards.Usecase
	log      *zap.Logger
}

func RegisterHandlers(mux *mux.Router, uc pWorkspaces.Usecase, boardsUC pBoards.Usecase, log *zap.Logger,
	checkAuth mw.Middleware, metrics mw.Middleware) {
	del := delivery{
		uc:       uc,
		boardsUC: boardsUC,
		log:      log,
	}

	const (
		workspacesPrefix = "/workspaces"
		workspacesPath   = constants.ApiPrefix + workspacesPrefix
		workspacePath    = workspacesPath + "/{id}"
	)

	mux.HandleFunc(workspacesPath, metrics(checkAuth(del.create))).Methods(http.MethodPost)
	mux.HandleFunc(workspacesPath, metrics(checkAuth(del.list))).Methods(http.MethodGet)

	mux.HandleFunc(workspacePath, metrics(checkAuth(del.get))).Methods(http.MethodGet)
	mux.HandleFunc(workspacePath, metrics(checkAuth(del.partialUpdate))).Methods(http.MethodPatch)
	mux.HandleFunc(workspacePath, metrics(checkAuth(del.delete))).Methods(http.MethodDelete)
}

// create godoc
//
//	@Summary		Create a new workspace
//	@Description	Create a new workspace
//	@Tags			workspaces
//	@Accept			json
//	@Produce		json
//	@Param			WorkspaceCreateData	body		createRequest	true	"Workspace create data"
//	@Success		200					{object}	createResponse	"Created workspace data."
//	@Failure		400					{object}	http.JSONError
//	@Failure		401					{object}	http.JSONError
//	@Failure		405
//	@Failure		500
//	@Router			/workspaces [post]
//
//	@Security		cookieAuth
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

// list godoc
//
//	@Summary		Returns all workspaces with boards of current user
//	@Description	Returns all workspaces with boards of current user
//	@Tags			workspaces
//	@Produce		json
//	@Success		200	{object}	listResponse	"Workspaces data"
//	@Failure		400	{object}	http.JSONError
//	@Failure		401	{object}	http.JSONError
//	@Failure		405
//	@Failure		500
//	@Router			/workspaces [get]
//
//	@Security		cookieAuth
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

	response := listResponse{Workspaces: make([]workspaceResponse, len(workspaces))}

	for i, _ := range workspaces {
		response.Workspaces[i].ID = workspaces[i].ID
		response.Workspaces[i].Title = workspaces[i].Title
		response.Workspaces[i].Description = workspaces[i].Description
		response.Workspaces[i].CreatedAt = workspaces[i].CreatedAt
		response.Workspaces[i].UpdatedAt = workspaces[i].UpdatedAt

		boards, err := del.boardsUC.ListByWorkspace(r.Context(), workspaces[i].ID)
		if err != nil {
			pHTTP.HandleError(w, r, err)
			return
		}
		response.Workspaces[i].Boards = boards
	}

	//response := newListResponse(workspaces)
	pHTTP.SendJSON(w, r, http.StatusOK, response)
}

// get godoc
//
//	@Summary		Returns workspace by id
//	@Description	Returns workspace by id
//	@Tags			workspaces
//	@Produce		json
//	@Param			id	path		int			true	"Workspace ID"
//	@Success		200	{object}	getResponse	"Workspace data"
//	@Failure		400	{object}	http.JSONError
//	@Failure		401	{object}	http.JSONError
//	@Failure		404	{object}	http.JSONError
//	@Failure		405
//	@Failure		500
//	@Router			/workspaces/{id} [get]
//
//	@Security		cookieAuth
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

// partialUpdate godoc
//
//	@Summary		Partial update of workspace
//	@Description	Partial update of workspace
//	@Tags			workspaces
//	@Accept			json
//	@Produce		json
//	@Param			id					path		int						true	"Workspace ID"
//	@Param			WorkspaceUpdateData	body		partialUpdateRequest	true	"Workspace data to update"
//	@Success		200					{object}	getResponse				"Updated workspace data."
//	@Failure		400					{object}	http.JSONError
//	@Failure		401					{object}	http.JSONError
//	@Failure		405
//	@Failure		500
//	@Router			/workspaces/{id}  [patch]
//
//	@Security		cookieAuth
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

// delete godoc
//
//	@Summary		Delete workspace by id
//	@Description	Delete workspace by id
//	@Tags			workspaces
//	@Produce		json
//	@Param			id	path	int	true	"Workspace ID"
//	@Success		204	"Workspace deleted successfully"
//	@Failure		400	{object}	http.JSONError
//	@Failure		401	{object}	http.JSONError
//	@Failure		404	{object}	http.JSONError
//	@Failure		405
//	@Failure		500
//	@Router			/workspaces/{id} [delete]
//
//	@Security		cookieAuth
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
