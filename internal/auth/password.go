package auth

import (
	"crypto/sha256"
	"fmt"
)

func Hash(password string) string {
	hash := sha256.Sum256([]byte(password))

	return fmt.Sprintf("%x", hash[:])
}

func CheckPassword(password, hash string) bool {
	passwordHash := sha256.Sum256([]byte(password))

	return fmt.Sprintf("%x", passwordHash[:]) == hash
}
