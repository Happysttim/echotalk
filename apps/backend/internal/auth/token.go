package auth

import (
	"time"

	"echotalk/internal/config"
	"echotalk/internal/errors"
	"echotalk/internal/model"

	"github.com/golang-jwt/jwt/v5"
)

type JWTToken struct {
	AccessToken  string
	RefreshToken string
}

func CreateToken(userID string) (*JWTToken, error) {
	accessToken, accessTokenErr := CreateAccessToken(userID)
	refreshToken, refreshTokenErr := CreateRefreshToken(userID)

	if accessTokenErr != nil {
		return nil, errors.ErrCreateAccessToken
	}

	if refreshTokenErr != nil {
		return nil, errors.ErrCreateRefreshToken
	}

	return &JWTToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func CreateAccessToken(userID string) (string, error) {
	config := config.Config
	secretKey := []byte(config.SecretKey)

	if len(secretKey) == 0 {
		return "", errors.ErrConfigSecretKey
	}

	now := time.Now()
	accessClaims := &model.AccessTokenClaims{
		UserID:    userID,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "echotalk-auth",
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(now.Add(30 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		accessClaims,
	)

	tokenString, err := token.SignedString(token)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func CreateRefreshToken(userID string) (string, error) {
	config := config.Config
	secretKey := []byte(config.SecretKey)

	if len(secretKey) == 0 {
		return "", errors.ErrConfigSecretKey
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

	tokenString, err := token.SignedString(token)
	if err != nil {
		return "", err
	}

	return tokenString, nil
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
		return nil, errors.ErrInvalidToken
	}

	claims = token.Claims.(*model.AccessTokenClaims)
	if claims.TokenType != "access" {
		return nil, errors.ErrInvalidTokenType
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
		return nil, errors.ErrInvalidToken
	}

	claims = token.Claims.(*model.RefreshTokenClaims)
	if claims.TokenType != "refresh" {
		return nil, errors.ErrInvalidTokenType
	}
	return claims, nil
}
