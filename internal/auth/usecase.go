package auth

import "github.com/SlavaShagalov/my-trello-backend/internal/models"

type SignInParams struct {
	Username string
	Password string
}

type SignUpParams struct {
	Name     string
	Username string
	Email    string
	Password string
}

type Usecase interface {
	SignIn(params *SignInParams) (models.User, string, error)
	SignUp(params *SignUpParams) (models.User, string, error)
	CheckAuth(userID int, authToken string) (int, error)
	Logout(userID int, authToken string) error
}
