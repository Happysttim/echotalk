package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Survey struct {
	ID        bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title     string        `json:"title" bson:"title"`
	Content   string        `json:"content" bson:"content"`
	IsPublic  bool          `json:"is_public" bson:"is_public"`
	Closed    bool          `json:"closed" bson:"closed"`
	AuthorID  bson.ObjectID `json:"author_id" bson:"author_id"`
	ExpiresAt time.Time     `json:"expires_at" bson:"expires_at"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" bson:"updated_at"`
}

type CreateSurveyRequest struct {
	Title     string    `json:"title" binding:"required"`
	Content   string    `json:"content" binding:"required"`
	IsPublic  bool      `json:"is_public" binding:"required"`
	ExpiresAt time.Time `json:"expires_at" binding:"required"`
}

type UpdateSurveyRequest struct {
	SurveyID  string    `json:"survey_id" binding:"required"`
	Title     string    `json:"title" binding:"required"`
	Content   string    `json:"content" binding:"required"`
	IsPublic  bool      `json:"is_public" binding:"required"`
	Closed    bool      `json:"closed" binding:"required"`
	ExpiresAt time.Time `json:"expires_at" binding:"required"`
}
