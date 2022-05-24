package authentication

//go:generate mockgen -destination ../../mocks/hash_password/mocks.go -source=./hash_password.go -package=mockHashPassword

import (
	"crypto/sha256"
	"encoding/hex"
)

type HashPassword interface {
	HashPassword(pwd string) string
}

type sha256Adapter struct {
}

func NewSha256Hash() HashPassword {
	return &sha256Adapter{}
}

func (s *sha256Adapter) HashPassword(pwd string) string {
	hash := sha256.Sum256([]byte(pwd))
	hashedPwd := hex.EncodeToString(hash[:])
	return hashedPwd
}
