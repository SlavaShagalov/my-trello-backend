package http

import (
	pCards "github.com/SlavaShagalov/my-trello-backend/internal/cards"
	pLists "github.com/SlavaShagalov/my-trello-backend/internal/lists"
	mw "github.com/SlavaShagalov/my-trello-backend/internal/middleware"
	"github.com/SlavaShagalov/my-trello-backend/internal/pkg/constants"
	pErrors "github.com/SlavaShagalov/my-trello-backend/internal/pkg/errors"
	pHTTP "github.com/SlavaShagalov/my-trello-backend/internal/pkg/http"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type delivery struct {
	uc      pLists.Usecase
	cardsUC pCards.Usecase
	log     *zap.Logger
}

func RegisterHandlers(mux *mux.Router, uc pLists.Usecase, cardsUC pCards.Usecase, log *zap.Logger, checkAuth mw.Middleware, metrics mw.Middleware) {
	del := delivery{
		uc:      uc,
		cardsUC: cardsUC,
		log:     log,
	}

	const (
		boardListsPrefix = "/boards/{id}/lists"
		boardListsPath   = constants.ApiPrefix + boardListsPrefix

		listsPrefix = "/lists"
		listsPath   = constants.ApiPrefix + listsPrefix
		listPath    = listsPath + "/{id}"
	)

	mux.HandleFunc(boardListsPath, metrics(checkAuth(del.create))).Methods(http.MethodPost)
	mux.HandleFunc(boardListsPath, metrics(checkAuth(del.listByBoard))).Methods(http.MethodGet)

	mux.HandleFunc(listsPath, metrics(checkAuth(del.list))).Methods(http.MethodGet).
		Queries("title", "{title}")

	mux.HandleFunc(listPath, metrics(checkAuth(del.get))).Methods(http.MethodGet)
	mux.HandleFunc(listPath, metrics(checkAuth(del.partialUpdate))).Methods(http.MethodPatch)
	mux.HandleFunc(listPath, metrics(checkAuth(del.delete))).Methods(http.MethodDelete)
}

// create godoc
//
//	@Summary		Create a new list
//	@Description	Create a new list
//	@Tags			boards
//	@Accept			json
//	@Produce		json
//	@Param			id				path		int				true	"Board ID"
//	@Param			ListCreateData	body		createRequest	true	"List create data"
//	@Success		200				{object}	createResponse	"Created list data."
//	@Failure		400				{object}	http.JSONError
//	@Failure		401				{object}	http.JSONError
//	@Failure		405
//	@Failure		500
//	@Router			/boards/{id}/lists [post]
//
//	@Security		cookieAuth
func (del *delivery) create(w http.ResponseWriter, r *http.Request) {
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

	var request createRequest
	err = request.UnmarshalJSON(body)
	if err != nil {
		pHTTP.HandleError(w, r, pErrors.ErrReadBody)
		return
	}

	params := pLists.CreateParams{
		Title:   request.Title,
		BoardID: boardID,
	}

	list, err := del.uc.Create(&params)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	response := newCreateResponse(&list)
	pHTTP.SendJSON(w, r, http.StatusOK, response)
}

// listByWorkspace godoc
//
//	@Summary		Returns lists by board id
//	@Description	Returns lists by board id
//	@Tags			boards
//	@Produce		json
//	@Param			id	path		int				true	"Board ID"
//	@Success		200	{object}	listResponse	"Lists data"
//	@Failure		400	{object}	http.JSONError
//	@Failure		401	{object}	http.JSONError
//	@Failure		405
//	@Failure		500
//	@Router			/boards/{id}/lists [get]
//
//	@Security		cookieAuth
func (del *delivery) listByBoard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	boardID, err := strconv.Atoi(vars["id"])
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	lists, err := del.uc.ListByBoard(boardID)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	response := listResponse{Lists: make([]itemResponse, len(lists))}

	for i, _ := range lists {
		response.Lists[i].ID = lists[i].ID
		response.Lists[i].BoardID = lists[i].BoardID
		response.Lists[i].Title = lists[i].Title
		response.Lists[i].Position = lists[i].Position
		response.Lists[i].CreatedAt = lists[i].CreatedAt
		response.Lists[i].UpdatedAt = lists[i].UpdatedAt

		cards, err := del.cardsUC.ListByList(lists[i].ID)
		if err != nil {
			pHTTP.HandleError(w, r, err)
			return
		}
		response.Lists[i].Cards = cards
	}

	//response := newListResponse(lists)
	pHTTP.SendJSON(w, r, http.StatusOK, response)
}

// list godoc
//
//	@Summary		Returns lists by board id
//	@Description	Returns lists by board id
//	@Tags			lists
//	@Produce		json
//	@Param			title	query		string			true	"Title filter"
//	@Success		200		{object}	listResponse	"Lists data"
//	@Failure		400		{object}	http.JSONError
//	@Failure		401		{object}	http.JSONError
//	@Failure		405
//	@Failure		500
//	@Router			/lists [get]
//
//	@Security		cookieAuth
func (del *delivery) list(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(mw.ContextUserID).(int)
	if !ok {
		pHTTP.HandleError(w, r, pErrors.ErrReadBody)
		return
	}

	title := r.FormValue("title")

	lists, err := del.uc.ListByTitle(title, userID)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	response := newListResponse(lists)
	pHTTP.SendJSON(w, r, http.StatusOK, response)
}

// get godoc
//
//	@Summary		Returns list by id
//	@Description	Returns list by id
//	@Tags			lists
//	@Produce		json
//	@Param			id	path		int			true	"Board ID"
//	@Success		200	{object}	getResponse	"Board data"
//	@Failure		400	{object}	http.JSONError
//	@Failure		401	{object}	http.JSONError
//	@Failure		404	{object}	http.JSONError
//	@Failure		405
//	@Failure		500
//	@Router			/lists/{id} [get]
//
//	@Security		cookieAuth
func (del *delivery) get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	listID, err := strconv.Atoi(vars["id"])
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	list, err := del.uc.Get(listID)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	response := newGetResponse(&list)
	pHTTP.SendJSON(w, r, http.StatusOK, response)
}

// partialUpdate godoc
//
//	@Summary		Partial update of list
//	@Description	Partial update of list
//	@Tags			lists
//	@Accept			json
//	@Produce		json
//	@Param			id				path		int						true	"List ID"
//	@Param			ListUpdateData	body		partialUpdateRequest	true	"List data to update"
//	@Success		200				{object}	getResponse				"Updated list data."
//	@Failure		400				{object}	http.JSONError
//	@Failure		401				{object}	http.JSONError
//	@Failure		405
//	@Failure		500
//	@Router			/lists/{id}  [patch]
//
//	@Security		cookieAuth
func (del *delivery) partialUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	listID, err := strconv.Atoi(vars["id"])
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

	params := pLists.PartialUpdateParams{ID: listID}
	params.UpdateTitle = request.Title != nil
	if params.UpdateTitle {
		params.Title = *request.Title
	}
	params.UpdateBoardID = request.BoardID != nil
	if params.UpdateBoardID {
		params.BoardID = *request.BoardID
	}
	params.UpdatePosition = request.Position != nil
	if params.UpdatePosition {
		params.Position = *request.Position
	}

	list, err := del.uc.PartialUpdate(&params)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	response := newGetResponse(&list)
	pHTTP.SendJSON(w, r, http.StatusOK, response)
}

// delete godoc
//
//	@Summary		Delete list by id
//	@Description	Delete list by id
//	@Tags			lists
//	@Produce		json
//	@Param			id	path	int	true	"List ID"
//	@Success		204	"List deleted successfully"
//	@Failure		400	{object}	http.JSONError
//	@Failure		401	{object}	http.JSONError
//	@Failure		404	{object}	http.JSONError
//	@Failure		405
//	@Failure		500
//	@Router			/lists/{id} [delete]
//
//	@Security		cookieAuth
func (del *delivery) delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	listID, err := strconv.Atoi(vars["id"])
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	err = del.uc.Delete(listID)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
