package bcrypt

import (
	"context"
	pkgHasher "github.com/SlavaShagalov/my-trello-backend/internal/pkg/hasher"
	"github.com/SlavaShagalov/my-trello-backend/internal/pkg/opentel"
	"golang.org/x/crypto/bcrypt"
)

const (
	componentName = "BCrypt Hasher"
)

type hasher struct{}

func New() pkgHasher.Hasher {
	return &hasher{}
}

func (h *hasher) GetHashedPassword(ctx context.Context, password string) (string, error) {
	_, span := opentel.Tracer.Start(ctx, componentName+" "+"GetHashedPassword")
	defer span.End()

	pswd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(pswd), err
}

func (h *hasher) CompareHashAndPassword(ctx context.Context, hashedPassword, password string) error {
	_, span := opentel.Tracer.Start(ctx, componentName+" "+"CompareHashAndPassword")
	defer span.End()

	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
