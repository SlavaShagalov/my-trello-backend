package e2e

import (
	"context"
	cardsDelivery "github.com/SlavaShagalov/my-trello-backend/internal/cards/delivery/http"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"log"
	"net/http"
	"time"

	"github.com/ozontech/cute"
)

func (s *EndToEndSuite) TestEditCard(t provider.T) {
	updateTitle := "Golang"
	updateContent := "Doc: ...; Log: ...; Git: ..."
	updateReq := cardsDelivery.PartialUpdateRequest{
		Title:   &updateTitle,
		Content: &updateContent,
	}

	resetTitle := "Определение требований"
	resetContent := "Определить основные требования к проект"
	resetReq := cardsDelivery.PartialUpdateRequest{
		Title:   &resetTitle,
		Content: &resetContent,
	}

	cute.NewTestBuilder().
		CreateStep("Get workspaces").
		RequestBuilder(
			cute.WithURI("http://test-api:8000/api/v1/workspaces"),
			cute.WithMethod(http.MethodGet),
			cute.WithHeadersKV("Cookie", s.authCookie),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusOK).
		NextTest().
		CreateStep("Get boards by workspace ID").
		RequestBuilder(
			cute.WithURI("http://test-api:8000/api/v1/workspaces/1/boards"),
			cute.WithMethod(http.MethodGet),
			cute.WithHeadersKV("Cookie", s.authCookie),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusOK).
		NextTest().
		CreateStep("Get lists by board ID").
		RequestBuilder(
			cute.WithURI("http://test-api:8000/api/v1/boards/1/lists"),
			cute.WithMethod(http.MethodGet),
			cute.WithHeadersKV("Cookie", s.authCookie),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusOK).
		NextTest().
		CreateStep("Get cards by list ID").
		RequestBuilder(
			cute.WithURI("http://test-api:8000/api/v1/lists/1/cards"),
			cute.WithMethod(http.MethodGet),
			cute.WithHeadersKV("Cookie", s.authCookie),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusOK).
		NextTest().
		CreateStep("Update card").
		RequestBuilder(
			cute.WithURI("http://test-api:8000/api/v1/cards/1"),
			cute.WithMethod(http.MethodPatch),
			cute.WithHeadersKV("Cookie", s.authCookie),
			cute.WithMarshalBody(updateReq),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusOK).
		NextTest().
		CreateStep("Check update").
		RequestBuilder(
			cute.WithURI("http://test-api:8000/api/v1/cards/1"),
			cute.WithMethod(http.MethodGet),
			cute.WithHeadersKV("Cookie", s.authCookie),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusOK).
		NextTest().
		CreateStep("Reset card values").
		RequestBuilder(
			cute.WithURI("http://test-api:8000/api/v1/cards/1"),
			cute.WithMethod(http.MethodPatch),
			cute.WithHeadersKV("Cookie", s.authCookie),
			cute.WithMarshalBody(resetReq),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusOK).
		NextTest().
		CreateStep("Check reset").
		RequestBuilder(
			cute.WithURI("http://test-api:8000/api/v1/cards/1"),
			cute.WithMethod(http.MethodGet),
			cute.WithHeadersKV("Cookie", s.authCookie),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusOK).
		ExecuteTest(context.Background(), t)

	log.Println(s.authCookie)
}
