package hasher

type Hasher interface {
	GetHashedPassword(password string) (string, error)
	CompareHashAndPassword(hashedPassword, password string) error
}
