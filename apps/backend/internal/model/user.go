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

type User struct {
	ID           bson.ObjectID `bson:"_id,omitempty"`
	AuthProvider AuthProvider  `bson:"auth_provider"`
	ProviderID   string        `bson:"provider_id"`
	Email        string        `bson:"email"`
	PasswordHash string        `bson:"password_hash"`
	CreatedAt    time.Time     `bson:"created_at"`
	UpdatedAt    time.Time     `bson:"updated_at"`
}

type CreateUserRequest struct {
	AuthProvider AuthProvider `json:"auth_provider" binding:"required"`
	Email        string       `json:"email" binding:"required,email"`
	Password     string       `json:"password"`
	Code         string       `json:"code"`
}
