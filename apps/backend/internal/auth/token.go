package auth

import (
	"errors"
	"time"

	"echotalk/internal/config"
	"echotalk/internal/model"

	"github.com/golang-jwt/jwt/v5"
)

func CreateAccessToken(userID string) (string, error) {
	config := config.Config
	secretKey := []byte(config.SecretKey)

	if len(secretKey) == 0 {
		return "", errors.New("secret key is not set in the configuration")
	}

	now := time.Now()
	accessClaims := &model.AccessTokenClaims{
		UserID:    userID,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "echotalk-auth",
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		accessClaims,
	)

	return token.SignedString(token)
}

func CreateRefreshToken(userID string) (string, error) {
	config := config.Config
	secretKey := []byte(config.SecretKey)

	if len(secretKey) == 0 {
		return "", errors.New("secret key is not set in the configuration")
	}

	now := time.Now()
	refreshClaims := &model.RefreshTokenClaims{
		UserID:    userID,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "echotalk-auth",
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(now.Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		refreshClaims,
	)

	return token.SignedString(secretKey)
}

func ParseAccessToken(tokenString string) (*model.AccessTokenClaims, error) {
	var claims *model.AccessTokenClaims

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Config.SecretKey), nil
		},
	)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims = token.Claims.(*model.AccessTokenClaims)
	if claims.TokenType != "access" {
		return nil, errors.New("invalid token type")
	}
	return claims, nil
}

func ParseRefreshToken(tokenString string) (*model.RefreshTokenClaims, error) {
	var claims *model.RefreshTokenClaims

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Config.SecretKey), nil
		},
	)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims = token.Claims.(*model.RefreshTokenClaims)
	if claims.TokenType != "refresh" {
		return nil, errors.New("invalid token type")
	}
	return claims, nil
}
