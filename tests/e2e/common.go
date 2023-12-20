package e2e

import (
	"context"
	authDelivery "github.com/SlavaShagalov/my-trello-backend/internal/auth/delivery/http"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/ozontech/cute"
	"log"
	"net/http"
	"time"
)

type EndToEndSuite struct {
	suite.Suite

	authCookie string
}

func (s *EndToEndSuite) BeforeAll(t provider.T) {
	cute.NewTestBuilder().
		Title("Sign In").
		Tags("card", "edit").
		Create().
		RequestBuilder(
			cute.WithURI("http://test-api:8000/api/v1/auth/signin"),
			cute.WithMethod(http.MethodPost),
			cute.WithMarshalBody(authDelivery.SignInRequest{Username: "slava", Password: "12345678"}),
		).
		After(func(response *http.Response, errors []error) error {
			respHeader := response.Header
			s.authCookie = respHeader.Get("Set-Cookie")
			return nil
		}).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusOK).
		ExecuteTest(context.Background(), t)

	log.Println(s.authCookie)
}
