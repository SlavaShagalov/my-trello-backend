package bcrypt

import (
	pkgHasher "github.com/SlavaShagalov/my-trello-backend/internal/pkg/hasher"
	"golang.org/x/crypto/bcrypt"
)

type hasher struct{}

func NewHasher() pkgHasher.Hasher {
	return &hasher{}
}

func (h *hasher) GetHashedPassword(password string) (string, error) {
	pswd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(pswd), err
}

func (h *hasher) CompareHashAndPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
