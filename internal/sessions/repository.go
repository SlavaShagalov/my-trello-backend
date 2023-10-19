package sessions

type Repository interface {
	Create(userID int) (string, error)
	Get(userID int, authToken string) (int, error)
	Delete(userID int, authToken string) error
}
