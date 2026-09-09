package service

import (
	"context"
	"echotalk/internal/errors"
	"echotalk/internal/model"
	"echotalk/internal/repositories"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type SurveyKeyword string

type SurveyService struct {
	surveyRepo  repositories.SurveyRepository
	userService *UserService
}

func NewSurveyService(surveyRepo repositories.SurveyRepository, userService *UserService) *SurveyService {
	return &SurveyService{
		surveyRepo:  surveyRepo,
		userService: userService,
	}
}

func (s *SurveyService) CreateSurvey(ctx context.Context, payload *model.CreateSurveyRequest, authorID string) (*model.Survey, error) {
	if payload == nil {
		return nil, errors.ErrInvalidInput
	}

	user, err := s.userService.GetUserByID(ctx, authorID)

	if err != nil {
		return nil, errors.ErrInternalServer
	}

	if user == nil {
		return nil, errors.ErrInvalidUser
	}

	doc := &model.Survey{
		Title:     payload.Title,
		Content:   payload.Content,
		AuthorID:  user.ID,
		IsPublic:  payload.IsPublic,
		Closed:    false,
		ExpiresAt: payload.ExpiresAt,
	}
	return s.surveyRepo.Create(ctx, doc)
}

func (s *SurveyService) GetSurveyByID(ctx context.Context, surveyID string) (*model.Survey, error) {
	if surveyID == "" {
		return nil, errors.ErrInvalidInput
	}
	survey, err := s.surveyRepo.FindByID(ctx, surveyID)
	if err != nil {
		return nil, err
	}

	if survey == nil {
		return nil, errors.ErrSurveyNotFound
	}
	return survey, nil
}

func (s *SurveyService) GetFilteredSurveys(ctx context.Context, filter bson.M, opts ...options.Lister[options.FindOptions]) ([]*model.Survey, error) {
	if filter == nil {
		return nil, errors.ErrInvalidInput
	}

	var m bson.M
	marshal, err := bson.Marshal(filter)

	if err != nil {
		return nil, errors.ErrInvalidInput
	}

	if err := bson.Unmarshal(marshal, &m); err != nil {
		return nil, errors.ErrInvalidInput
	}

	return s.surveyRepo.FindByFilter(ctx, m, opts...)
}

func (s *SurveyService) GetAllSurveys(ctx context.Context) ([]*model.Survey, error) {
	return s.surveyRepo.FindAll(ctx)
}

func (s *SurveyService) UpdateSurvey(ctx context.Context, payload *model.UpdateSurveyRequest) error {
	if payload == nil {
		return errors.ErrInvalidInput
	}

	survey, err := s.surveyRepo.FindByID(ctx, payload.SurveyID)
	if err != nil {
		return err
	}

	if survey == nil {
		return errors.ErrSurveyNotFound
	}

	survey.Title = payload.Title
	survey.Content = payload.Content
	survey.IsPublic = *payload.IsPublic
	survey.ExpiresAt = payload.ExpiresAt
	survey.Closed = *payload.Closed

	survey.UpdatedAt = time.Now()

	return s.surveyRepo.Update(ctx, survey)
}

func (s *SurveyService) DeleteSurvey(ctx context.Context, surveyID string) error {
	if surveyID == "" {
		return errors.ErrInvalidInput
	}

	return s.surveyRepo.Delete(ctx, surveyID)
}

func (s *SurveyService) CheckExpire(ctx context.Context) error {
	return s.surveyRepo.CheckExpires(ctx)
}
