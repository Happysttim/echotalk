package service

import (
	"context"
	"echotalk/internal/errors"
	"echotalk/internal/model"
	"echotalk/internal/repositories"
	"time"
)

type AnswerService struct {
	answerRepo *repositories.MongoAnswerRepository
}

func NewAnswerService(answerRepo *repositories.MongoAnswerRepository) *AnswerService {
	return &AnswerService{
		answerRepo: answerRepo,
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

	if payload.Content == "" {
		return errors.ErrInvalidInput
	}

	if payload.AnswerID == "" {
		return errors.ErrAnswerNotFound
	}

	if payload.SurveyID == "" {
		return errors.ErrSurveyNotFound
	}

	answer, err := s.answerRepo.FindByID(ctx, payload.AnswerID)
	if err != nil {
		return err
	}

	answer.Content = payload.Content
	answer.UpdatedAt = time.Now()

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
