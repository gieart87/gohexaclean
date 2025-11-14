package service

// HashService defines the outbound port for password hashing
type HashService interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}
