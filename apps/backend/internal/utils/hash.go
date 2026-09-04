package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func Hash(target string) (string, error) {
	if target == "" {
		return "", errors.New("hash target cannot be empty")
	}

	hasher := sha256.New()
	_, err := hasher.Write([]byte(target))
	if err != nil {
		return "", err
	}

	shaResult := hex.EncodeToString(hasher.Sum(nil))
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(shaResult), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func Sha256Hash(target string) (string, error) {
	if target == "" {
		return "", errors.New("hash target cannot be empty")
	}

	hasher := sha256.New()
	_, err := hasher.Write([]byte(target))
	if err != nil {
		return "", err
	}

	shaResult := hex.EncodeToString(hasher.Sum(nil))
	return shaResult, nil
}

func CompareHash(src string, dst string) error {
	hasher := sha256.New()
	_, err := hasher.Write([]byte(dst))
	if err != nil {
		return err
	}
	shaResult := hex.EncodeToString(hasher.Sum(nil))
	return bcrypt.CompareHashAndPassword([]byte(src), []byte(shaResult))
}
