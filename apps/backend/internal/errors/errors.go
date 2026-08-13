package errors

import "errors"

var (
	ErrSurveyNotFound = errors.New("survey not found")
	ErrAnswerNotFound = errors.New("answer not found")
)

var (
	ErrInvalidSurvey      = errors.New("invalid survey data")
	ErrInvalidAnswer      = errors.New("invalid answer data")
	ErrInvalidUser        = errors.New("invalid user data")
	ErrInvalidInput       = errors.New("invalid input data")
	ErrInvalidProvider    = errors.New("invalid auth provider")
	ErrInvalidOAuthClaims = errors.New("invalid oauth claims")
	ErrInvalidIDToken     = errors.New("invalid id token")
	ErrInvalidToken       = errors.New("invalid token")
	ErrInvalidTokenType   = errors.New("invalid token type")
)

var (
	ErrExpiredToken = errors.New("expired token error")
)

var (
	ErrBadRequest = errors.New("bad request")
)

var (
	ErrCreateAccessToken  = errors.New("failed to create access token")
	ErrCreateRefreshToken = errors.New("failed to create refresh token")
)

var (
	ErrConfigSecretKey = errors.New("secret key is not set in the configuration")
)

var (
	ErrInternalServer = errors.New("internal server error")
)
