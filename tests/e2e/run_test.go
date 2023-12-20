package e2e

import (
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"testing"
)

func TestEndToEndSuite(t *testing.T) {
	suite.RunSuite(t, new(EndToEndSuite))
}
