package model

import (
	"echotalk/internal/errors"
	"echotalk/internal/utils"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Answer struct {
	ID             bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	SurveyID       bson.ObjectID `json:"survey_id" bson:"survey_id"`
	AuthorID       bson.ObjectID `json:"author_id" bson:"author_id"`
	IsAnonymous    bool          `json:"is_anonymous" bson:"is_anonymous"`
	Anonymous      string        `json:"anonymous" bson:"anonymous"`
	AnswerPassword string        `json:"-" bson:"answer_password"`
	Content        string        `json:"content" bson:"content"`
	RateUp         int           `json:"rate_up" bson:"rate_up"`
	CreatedAt      time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at" bson:"updated_at"`
}

type AnswerDetail struct {
	ID          bson.ObjectID `json:"id" bson:"_id"`
	SurveyID    bson.ObjectID `json:"survey_id" bson:"survey_id"`
	AuthorID    bson.ObjectID `json:"author_id" bson:"author_id"`
	IsAnonymous bool          `json:"is_anonymous" bson:"is_anonymous"`
	Anonymous   string        `json:"anonymous" bson:"anonymous"`
	Content     string        `json:"content" bson:"content"`
	RateUp      int           `json:"rate_up" bson:"rate_up"`
	CreatedAt   time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at" bson:"updated_at"`

	Author UserSummary `json:"author" bson:"author"`
}

type MyAnswerDetail struct {
	Survey struct {
		ID        bson.ObjectID  `json:"id,omitempty" bson:"_id,omitempty"`
		Title     string         `json:"title" bson:"title"`
		IsPublic  bool           `json:"is_public" bson:"is_public"`
		Closed    bool           `json:"closed" bson:"closed"`
		AuthorID  bson.ObjectID  `json:"author_id" bson:"author_id"`
		ExpiresAt time.Time      `json:"expires_at" bson:"expires_at"`
		CreatedAt time.Time      `json:"created_at" bson:"created_at"`
		UpdatedAt time.Time      `json:"updated_at" bson:"updated_at"`
		Answers   []AnswerDetail `json:"answers" bson:"answers"`
	} `json:"survey_answers" bson:"survey_answers"`
}

type CreateAnswerRequest struct {
	SurveyID       string `json:"survey_id" binding:"required"`
	Content        string `json:"content" binding:"required"`
	IsAnonymous    bool   `json:"is_anonymous"`
	Anonymous      string `json:"anonymous"`
	AnswerPassword string `json:"answer_password"`
	AuthorID       string `json:"-"`
}

type UpdateAnswerRequest struct {
	SurveyID string `json:"survey_id" binding:"required"`
	AnswerID string `json:"answer_id" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

type DeleteAnswerRequest struct {
	AnswerID       string `json:"answer_id" binding:"required"`
	AnswerPassword string `json:"answer_password"`
}

func (payload *CreateAnswerRequest) AnonymousName() string {
	if payload.IsAnonymous {
		if payload.Anonymous == "" {
			return "Anonymous"
		}
		return payload.Anonymous
	}

	return ""
}

func (payload *CreateAnswerRequest) AnonymousPassword() (string, error) {
	if payload.IsAnonymous {
		if payload.AnswerPassword == "" {
			return "", errors.ErrInvalidInput
		}
		sha256Pass, err := utils.Sha256Hash(payload.AnswerPassword)
		if err != nil {
			return "", err
		}

		return sha256Pass, nil
	}

	return "", nil
}
