package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Answer struct {
	ID        bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	SurveyID  bson.ObjectID `json:"survey_id" bson:"survey_id"`
	AuthorID  bson.ObjectID `json:"author_id" bson:"author_id"`
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

type DeleteAnswerRequest struct {
	AnswerID string `json:"answer_id" binding:"required"`
}
