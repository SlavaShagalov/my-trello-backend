package http

import (
	pCards "github.com/SlavaShagalov/my-trello-backend/internal/cards"
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
	uc  pCards.Usecase
	log *zap.Logger
}

func RegisterHandlers(mux *mux.Router, uc pCards.Usecase, log *zap.Logger, checkAuth mw.Middleware) {
	del := delivery{
		uc:  uc,
		log: log,
	}

	const (
		listCardsPrefix = "/lists/{id}/cards"
		listCardsPath   = constants.ApiPrefix + listCardsPrefix

		cardsPrefix = "/cards"
		cardsPath   = constants.ApiPrefix + cardsPrefix
		cardPath    = cardsPath + "/{id}"
	)

	mux.HandleFunc(listCardsPath, checkAuth(del.create)).Methods(http.MethodPost)
	mux.HandleFunc(listCardsPath, checkAuth(del.listByList)).Methods(http.MethodGet)

	mux.HandleFunc(cardsPath, checkAuth(del.list)).Methods(http.MethodGet).
		Queries("title", "{title}")

	mux.HandleFunc(cardPath, checkAuth(del.get)).Methods(http.MethodGet)
	mux.HandleFunc(cardPath, checkAuth(del.partialUpdate)).Methods(http.MethodPatch)
	mux.HandleFunc(cardPath, checkAuth(del.delete)).Methods(http.MethodDelete)
}

func (del *delivery) create(w http.ResponseWriter, r *http.Request) {
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

	var request createRequest
	err = request.UnmarshalJSON(body)
	if err != nil {
		pHTTP.HandleError(w, r, pErrors.ErrReadBody)
		return
	}

	params := pCards.CreateParams{
		Title:   request.Title,
		Content: request.Content,
		ListID:  listID,
	}

	card, err := del.uc.Create(&params)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	response := newCreateResponse(&card)
	pHTTP.SendJSON(w, r, http.StatusOK, response)
}

func (del *delivery) listByList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	listID, err := strconv.Atoi(vars["id"])
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	cards, err := del.uc.ListByList(listID)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	response := newCardResponse(cards)
	pHTTP.SendJSON(w, r, http.StatusOK, response)
}

func (del *delivery) list(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(mw.ContextUserID).(int)
	if !ok {
		pHTTP.HandleError(w, r, pErrors.ErrReadBody)
		return
	}

	title := r.FormValue("title")

	cards, err := del.uc.ListByTitle(title, userID)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	response := newCardResponse(cards)
	pHTTP.SendJSON(w, r, http.StatusOK, response)
}

func (del *delivery) get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cardID, err := strconv.Atoi(vars["id"])
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	card, err := del.uc.Get(cardID)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	response := newGetResponse(&card)
	pHTTP.SendJSON(w, r, http.StatusOK, response)
}

func (del *delivery) partialUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cardID, err := strconv.Atoi(vars["id"])
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

	params := pCards.PartialUpdateParams{ID: cardID}
	params.UpdateTitle = request.Title != nil
	if params.UpdateTitle {
		params.Title = *request.Title
	}
	params.UpdateContent = request.Content != nil
	if params.UpdateContent {
		params.Content = *request.Content
	}
	params.UpdateListID = request.ListID != nil
	if params.UpdateListID {
		params.ListID = *request.ListID
	}
	params.UpdatePosition = request.Position != nil
	if params.UpdatePosition {
		params.Position = *request.Position
	}

	card, err := del.uc.PartialUpdate(&params)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	response := newGetResponse(&card)
	pHTTP.SendJSON(w, r, http.StatusOK, response)
}

func (del *delivery) delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cardID, err := strconv.Atoi(vars["id"])
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	err = del.uc.Delete(cardID)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
