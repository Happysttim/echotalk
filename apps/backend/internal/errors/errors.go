package errors

import "errors"

var (
	ErrInvalidSurvey   = errors.New("invalid survey data")
	ErrSurveyNotFound  = errors.New("survey not found")
	ErrAnswerNotFound  = errors.New("answer not found")
	ErrInvalidAnswer   = errors.New("invalid answer data")
	ErrInvalidUser     = errors.New("invalid user data")
	ErrInvalidInput    = errors.New("invalid input data")
	ErrInvalidProvider = errors.New("invalid auth provider")
)
