package sessions

import "context"

type Repository interface {
	Create(ctx context.Context, userID int) (string, error)
	Get(ctx context.Context, userID int, authToken string) (int, error)
	Delete(ctx context.Context, userID int, authToken string) error
}
