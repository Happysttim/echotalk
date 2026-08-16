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
	answerRepo    *repositories.MongoAnswerRepository
	rateUpRepo    *repositories.MongoRateUpRepository
	userService   *UserService
	surveyService *SurveyService
}

func NewAnswerService(answerRepo *repositories.MongoAnswerRepository, rateUpRepo *repositories.MongoRateUpRepository, userService *UserService, surveyService *SurveyService) *AnswerService {
	return &AnswerService{
		answerRepo:    answerRepo,
		rateUpRepo:    rateUpRepo,
		userService:   userService,
		surveyService: surveyService,
	}
}

func (s *AnswerService) CreateAnswer(ctx context.Context, payload *model.CreateAnswerRequest, authorID string) (*model.Answer, error) {
	if payload == nil {
		return nil, errors.ErrInvalidAnswer
	}

	user, err := s.userService.GetUserByID(ctx, authorID)

	if err != nil {
		return nil, errors.ErrInternalServer
	}

	if user == nil {
		return nil, errors.ErrInvalidUser
	}

	survey, err := s.surveyService.GetSurveyByID(ctx, payload.SurveyID)

	if err != nil {
		return nil, errors.ErrInternalServer
	}

	if survey == nil {
		return nil, errors.ErrInvalidSurvey
	}

	doc := &model.Answer{
		Content:   payload.Content,
		SurveyID:  survey.ID,
		AuthorID:  user.ID,
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

func (s *AnswerService) GetAnswersBySurvey(ctx context.Context, surveyID string) ([]*model.Answer, error) {
	if surveyID == "" {
		return nil, errors.ErrInvalidAnswer
	}

	answers, err := s.answerRepo.FindBySurveyID(ctx, surveyID)
	if err != nil {
		return nil, err
	}

	return answers, nil
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

	return s.answerRepo.Delete(ctx, answerID)
}

func (s *AnswerService) DeleteAnswerBySurvey(ctx context.Context, surveyID string) error {
	if surveyID == "" {
		return errors.ErrInvalidInput
	}

	return s.answerRepo.DeleteBySurveyID(ctx, surveyID)
}
