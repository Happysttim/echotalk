package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type AuthProvider string

const (
	AuthProviderLocal  AuthProvider = "local"
	AuthProviderGoogle AuthProvider = "google"
)

func (provider AuthProvider) IsValid() bool {
	switch provider {
	case AuthProviderGoogle, AuthProviderLocal:
		return true
	default:
		return false
	}
}

type User struct {
	ID           bson.ObjectID `bson:"_id,omitempty"`
	AuthProvider AuthProvider  `bson:"auth_provider"`
	ProviderID   string        `bson:"provider_id"`
	Email        string        `bson:"email"`
	PasswordHash string        `bson:"password_hash"`
	CreatedAt    time.Time     `bson:"created_at"`
	UpdatedAt    time.Time     `bson:"updated_at"`
}

type LocalAuthRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type GoogleAuthRequest struct {
	Code string `form:"code" binding:"required"`
}

type CreateUserCommand struct {
	AuthProvider AuthProvider
	ProviderID   string
	Email        string
	PasswordHash string
}
