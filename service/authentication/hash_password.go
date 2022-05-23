package authentication

import (
	"crypto/sha256"
	"encoding/hex"
)

type HashPassword interface {
	hashPassword(pwd string) string
}

type sha256Adapter struct {
}

func NewSha256Hash() HashPassword {
	return &sha256Adapter{}
}

func (s *sha256Adapter) hashPassword(pwd string) string {
	hash := sha256.Sum256([]byte(pwd))
	hashedPwd := hex.EncodeToString(hash[:])
	return hashedPwd
}
