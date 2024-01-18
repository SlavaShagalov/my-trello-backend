package bcrypt

import (
	"context"
	pkgHasher "github.com/SlavaShagalov/my-trello-backend/internal/pkg/hasher"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/bcrypt"
)

const (
	componentName = "BCrypt Hasher"
)

type hasher struct {
	tracer trace.Tracer
}

func New(tracer trace.Tracer) pkgHasher.Hasher {
	return &hasher{
		tracer: tracer,
	}
}

func (h *hasher) GetHashedPassword(ctx context.Context, password string) (string, error) {
	_, span := h.tracer.Start(ctx, componentName+" "+"GetHashedPassword")
	defer span.End()

	pswd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(pswd), err
}

func (h *hasher) CompareHashAndPassword(ctx context.Context, hashedPassword, password string) error {
	_, span := h.tracer.Start(ctx, componentName+" "+"CompareHashAndPassword")
	defer span.End()

	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
