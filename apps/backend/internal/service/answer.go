package service

import (
	"context"
	"echotalk/internal/errors"
	"echotalk/internal/model"
	"echotalk/internal/repositories"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type AnswerService struct {
	answerRepo *repositories.MongoAnswerRepository
	rateUpRepo *repositories.MongoRateUpRepository
}

func NewAnswerService(answerRepo *repositories.MongoAnswerRepository, rateUpRepo *repositories.MongoRateUpRepository) *AnswerService {
	return &AnswerService{
		answerRepo: answerRepo,
		rateUpRepo: rateUpRepo,
	}
}

func (s *AnswerService) CreateAnswer(ctx context.Context, payload *model.CreateAnswerRequest, author *model.User) (*model.Answer, error) {
	if payload == nil {
		return nil, errors.ErrInvalidAnswer
	}

	if payload.Content == "" {
		return nil, errors.ErrInvalidInput
	}

	if payload.SurveyID == "" {
		return nil, errors.ErrSurveyNotFound
	}

	if author == nil || author.ID.Hex() == "" {
		return nil, errors.ErrInvalidUser
	}

	doc := &model.Answer{
		Content:   payload.Content,
		Author:    author,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return s.answerRepo.Create(ctx, doc)
}

func (s *AnswerService) GetAnswerByID(ctx context.Context, answerID string) (*model.Answer, error) {
	if answerID == "" {
		return nil, errors.ErrInvalidInput
	}

	answer, err := s.answerRepo.FindByID(ctx, answerID)
	if err != nil {
		return nil, err
	}
	return answer, nil
}

func (s *AnswerService) UpdateAnswer(ctx context.Context, payload *model.UpdateAnswerRequest) error {
	if payload == nil {
		return errors.ErrInvalidAnswer
	}

	answer, err := s.answerRepo.FindByID(ctx, payload.AnswerID)
	if err != nil {
		return err
	}

	if answer == nil {
		return errors.ErrInvalidAnswer
	}

	answer.Content = payload.Content
	answer.UpdatedAt = time.Now()

	return s.answerRepo.Update(ctx, answer)
}

func (s *AnswerService) RateUp(ctx context.Context, answerID string, userID bson.ObjectID) error {
	if answerID == "" {
		return errors.ErrInvalidInput
	}

	answer, err := s.answerRepo.FindByID(ctx, answerID)
	if err != nil {
		return err
	}

	if answer == nil {
		return errors.ErrInvalidAnswer
	}

	rateUp := model.RateUp{
		AnswerID: answer.ID,
		UserID:   userID,
	}

	if err := s.rateUpRepo.RateUp(ctx, &rateUp); err != nil {
		return err
	}

	answer.RateUp += 1
	return s.answerRepo.Update(ctx, answer)
}

func (s *AnswerService) DeleteAnswer(ctx context.Context, answerID string) error {
	if answerID == "" {
		return errors.ErrInvalidInput
	}

	answer, err := s.answerRepo.FindByID(ctx, answerID)
	if err != nil {
		return err
	}

	if answer == nil {
		return errors.ErrAnswerNotFound
	}

	return s.answerRepo.Delete(ctx, answerID)
}
