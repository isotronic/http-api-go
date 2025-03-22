package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashByte, err := bcrypt.GenerateFromPassword([]byte(password), 13)
	if err != nil {
		return "", err
	}
	return string(hashByte), nil
}