package auth

import (
	"time"

	"echotalk/internal/config"
	"echotalk/internal/errors"
	"echotalk/internal/model"

	"github.com/golang-jwt/jwt/v5"
)

type JWTToken struct {
	AccessToken           string
	RefreshToken          string
	RefreshTokenExpiresIn time.Time
}

func CreateToken(userID string) (*JWTToken, error) {
	accessToken, accessTokenErr := CreateAccessToken(userID)
	refreshToken, expiresIn, refreshTokenErr := CreateRefreshToken(userID)

	if accessTokenErr != nil {
		return nil, errors.ErrCreateAccessToken
	}

	if refreshTokenErr != nil {
		return nil, errors.ErrCreateRefreshToken
	}

	return &JWTToken{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresIn: expiresIn,
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
			ExpiresAt: jwt.NewNumericDate(now.Add(5 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		accessClaims,
	)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func CreateRefreshToken(userID string) (string, time.Time, error) {
	config := config.Config
	secretKey := []byte(config.SecretKey)

	if len(secretKey) == 0 {
		return "", time.Time{}, errors.ErrConfigSecretKey
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

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, refreshClaims.ExpiresAt.Time, nil
}

func CreateVerifiedToken(email string) (string, error) {
	config := config.Config
	secretKey := []byte(config.SecretKey)

	if len(secretKey) == 0 {
		return "", errors.ErrConfigSecretKey
	}

	now := time.Now()
	verifiedClaims := &model.VerifiedTokenClaims{
		Email:     email,
		TokenType: "verified",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "echotalk-register-verified",
			Subject:   email,
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		verifiedClaims,
	)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseAccessToken(tokenString string) (*model.AccessTokenClaims, error) {
	claims := &model.AccessTokenClaims{}

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
	claims := &model.RefreshTokenClaims{}

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

func ParseVerifiedToken(tokenString string) (*model.VerifiedTokenClaims, error) {
	claims := &model.VerifiedTokenClaims{}

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

	claims = token.Claims.(*model.VerifiedTokenClaims)
	if claims.TokenType != "verified" {
		return nil, errors.ErrInvalidTokenType
	}
	return claims, nil
}
