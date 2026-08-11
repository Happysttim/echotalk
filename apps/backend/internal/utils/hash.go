package utils

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(target string) (string, error) {
	if target == "" {
		return "", errors.New("hash target cannot be empty")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(target), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}
