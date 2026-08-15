package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Answer struct {
	ID        bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Author    *User         `json:"author" bson:"author"`
	Content   string        `json:"content" bson:"content"`
	RateUp    int           `json:"rate_up" bson:"rate_up"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" bson:"updated_at"`
}

type CreateAnswerRequest struct {
	SurveyID string `json:"survey_id" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

type UpdateAnswerRequest struct {
	SurveyID string `json:"survey_id" binding:"required"`
	AnswerID string `json:"answer_id" binding:"required"`
	Content  string `json:"content" binding:"required"`
}
