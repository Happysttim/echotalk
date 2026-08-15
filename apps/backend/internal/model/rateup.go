package model

import "go.mongodb.org/mongo-driver/v2/bson"

type RateUp struct {
	ID       bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	AnswerID bson.ObjectID `json:"answer_id" bson:"answer_id"`
	UserID   bson.ObjectID `json:"user_ud" bson:"user_id"`
}

type RateUpRequest struct {
	AnswerID string `json:"answer_id" binding:"required"`
}
