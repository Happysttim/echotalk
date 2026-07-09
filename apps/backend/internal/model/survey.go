package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Answer struct {
	ID        bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Author    User          `json:"author" bson:"author"`
	Content   string        `json:"content" bson:"content"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" bson:"updated_at"`
}

type Survey struct {
	ID        bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title     string        `json:"title" bson:"title"`
	Content   string        `json:"content" bson:"content"`
	Answers   []Answer      `json:"answers" bson:"answers"`
	Author    User          `json:"author" bson:"author"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" bson:"updated_at"`
}

type CreateSurveyRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	Author  string `json:"author" binding:"required"`
}

type UpdateSurveyRequest struct {
	SurveyID string `json:"survey_id" binding:"required"`
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

type CreateAnswerRequest struct {
	SurveyID string `json:"survey_id" binding:"required"`
	Author   string `json:"author" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

type UpdateAnswerRequest struct {
	SurveyID string `json:"survey_id" binding:"required"`
	AnswerID string `json:"answer_id" binding:"required"`
	Content  string `json:"content" binding:"required"`
}
