package service

import (
	"context"
	"echotalk/internal/errors"
	"echotalk/internal/model"
	"echotalk/internal/repositories"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type SurveyKeyword string

const (
	SurveyKeywordID        SurveyKeyword = "_id"
	SurveyKeywordTitle     SurveyKeyword = "title"
	SurveyKeywordContent   SurveyKeyword = "content"
	SurveyKeywordAuthor    SurveyKeyword = "author"
	SurveyKeywordCreatedAt SurveyKeyword = "created_at"
)

type SurveyFilter map[SurveyKeyword]interface{}

type SurveyService struct {
	surveyRepo *repositories.MongoSurveyRepository
}

func NewSurveyService(surveyRepo *repositories.MongoSurveyRepository) *SurveyService {
	return &SurveyService{
		surveyRepo: surveyRepo,
	}
}

func (s *SurveyService) CreateSurvey(ctx context.Context, payload *model.CreateSurveyRequest, author *model.User) (*model.Survey, error) {
	if payload == nil {
		return nil, errors.ErrInvalidSurvey
	}

	if payload.Title == "" || payload.Content == "" {
		return nil, errors.ErrInvalidInput
	}

	if author == nil || author.ID.Hex() == "" {
		return nil, errors.ErrInvalidUser
	}

	doc := &model.Survey{
		Title:   payload.Title,
		Content: payload.Content,
		Author:  author,
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

func (s *SurveyService) GetFilteredSurveys(ctx context.Context, filter SurveyFilter) ([]*model.Survey, error) {
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

	return s.surveyRepo.FindByFilter(ctx, m)
}

func (s *SurveyService) GetAllSurveys(ctx context.Context) ([]*model.Survey, error) {
	return s.surveyRepo.FindAll(ctx)
}

func (s *SurveyService) UpdateSurvey(ctx context.Context, payload *model.UpdateSurveyRequest) error {
	if payload == nil {
		return errors.ErrInvalidSurvey
	}

	if payload.SurveyID == "" || payload.Title == "" || payload.Content == "" {
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
	survey.UpdatedAt = time.Now()

	return s.surveyRepo.Update(ctx, survey)
}

func (s *SurveyService) DeleteSurvey(ctx context.Context, surveyID string) error {
	if surveyID == "" {
		return errors.ErrInvalidInput
	}

	survey, err := s.surveyRepo.FindByID(ctx, surveyID)
	if err != nil {
		return err
	}

	if survey == nil {
		return errors.ErrSurveyNotFound
	}

	return s.surveyRepo.Delete(ctx, surveyID)
}
