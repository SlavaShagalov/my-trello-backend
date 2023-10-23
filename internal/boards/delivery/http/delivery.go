package http

import (
	"bytes"
	pBoards "github.com/SlavaShagalov/my-trello-backend/internal/boards"
	mw "github.com/SlavaShagalov/my-trello-backend/internal/middleware"
	"github.com/SlavaShagalov/my-trello-backend/internal/pkg/constants"
	pErrors "github.com/SlavaShagalov/my-trello-backend/internal/pkg/errors"
	pHTTP "github.com/SlavaShagalov/my-trello-backend/internal/pkg/http"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
)

type delivery struct {
	uc  pBoards.Usecase
	log *zap.Logger
}

func RegisterHandlers(mux *mux.Router, uc pBoards.Usecase, log *zap.Logger, checkAuth mw.Middleware) {
	del := delivery{
		uc:  uc,
		log: log,
	}

	const (
		workspaceBoardsPrefix = "/workspaces/{id}/boards"
		workspaceBoardsPath   = constants.ApiPrefix + workspaceBoardsPrefix

		boardsPrefix = "/boards"
		boardsPath   = constants.ApiPrefix + boardsPrefix
		boardPath    = boardsPath + "/{id}"

		backgroundPath = boardPath + "/background"
	)

	mux.HandleFunc(workspaceBoardsPath, checkAuth(del.create)).Methods(http.MethodPost)
	mux.HandleFunc(workspaceBoardsPath, checkAuth(del.listByWorkspace)).Methods(http.MethodGet)

	mux.HandleFunc(boardsPath, checkAuth(del.list)).Methods(http.MethodGet).
		Queries("title", "{title}")

	mux.HandleFunc(boardPath, checkAuth(del.get)).Methods(http.MethodGet)
	mux.HandleFunc(boardPath, checkAuth(del.partialUpdate)).Methods(http.MethodPatch)
	mux.HandleFunc(backgroundPath, checkAuth(del.updateBackground)).Methods(http.MethodPut)
	mux.HandleFunc(boardPath, checkAuth(del.delete)).Methods(http.MethodDelete)
}

func (del *delivery) create(w http.ResponseWriter, r *http.Request) {
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

	var request createRequest
	err = request.UnmarshalJSON(body)
	if err != nil {
		pHTTP.HandleError(w, r, pErrors.ErrReadBody)
		return
	}

	params := pBoards.CreateParams{
		Title:       request.Title,
		Description: request.Description,
		WorkspaceID: workspaceID,
	}

	board, err := del.uc.Create(&params)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	response := newCreateResponse(&board)
	pHTTP.SendJSON(w, r, http.StatusOK, response)
}

func (del *delivery) listByWorkspace(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	workspaceID, err := strconv.Atoi(vars["id"])
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	boards, err := del.uc.ListByWorkspace(workspaceID)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	response := newListResponse(boards)
	pHTTP.SendJSON(w, r, http.StatusOK, response)
}

func (del *delivery) list(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(mw.ContextUserID).(int)
	if !ok {
		pHTTP.HandleError(w, r, pErrors.ErrReadBody)
		return
	}

	title := r.FormValue("title")

	boards, err := del.uc.ListByTitle(title, userID)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	response := newListResponse(boards)
	pHTTP.SendJSON(w, r, http.StatusOK, response)
}

func (del *delivery) get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	boardID, err := strconv.Atoi(vars["id"])
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	board, err := del.uc.Get(boardID)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	response := newGetResponse(&board)
	pHTTP.SendJSON(w, r, http.StatusOK, response)
}

func (del *delivery) partialUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	boardID, err := strconv.Atoi(vars["id"])
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

	params := pBoards.PartialUpdateParams{ID: boardID}
	params.UpdateTitle = request.Title != nil
	if params.UpdateTitle {
		params.Title = *request.Title
	}
	params.UpdateDescription = request.Description != nil
	if params.UpdateDescription {
		params.Description = *request.Description
	}

	board, err := del.uc.PartialUpdate(&params)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	response := newGetResponse(&board)
	pHTTP.SendJSON(w, r, http.StatusOK, response)
}

func (del *delivery) updateBackground(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	file, header, err := r.FormFile("background")
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return

	}

	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, file)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	board, err := del.uc.UpdateBackground(userID, buf.Bytes(), header.Filename)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	response := newGetResponse(board)
	pHTTP.SendJSON(w, r, http.StatusOK, response)
}

func (del *delivery) delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	boardID, err := strconv.Atoi(vars["id"])
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	err = del.uc.Delete(boardID)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
