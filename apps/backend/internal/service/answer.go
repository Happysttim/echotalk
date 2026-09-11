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

type AnswerService struct {
	answerRepo  repositories.AnswerRepository
	rateUpRepo  repositories.RateUpRepository
	userService *UserService
}

func NewAnswerService(answerRepo repositories.AnswerRepository, rateUpRepo repositories.RateUpRepository, userService *UserService) *AnswerService {
	return &AnswerService{
		answerRepo:  answerRepo,
		rateUpRepo:  rateUpRepo,
		userService: userService,
	}
}

func (s *AnswerService) CreateAnswer(ctx context.Context, payload *model.CreateAnswerRequest) (*model.Answer, error) {
	if payload == nil {
		return nil, errors.ErrInvalidAnswer
	}

	var authorID bson.ObjectID

	if !payload.IsAnonymous {
		if payload.AuthorID == "" {
			return nil, errors.ErrInvalidUser
		}
		user, err := s.userService.GetUserByID(ctx, payload.AuthorID)

		if err != nil {
			return nil, errors.ErrInternalServer
		}

		if user == nil {
			return nil, errors.ErrInvalidUser
		}

		authorID = user.ID
	}

	surveyObjectID, err := bson.ObjectIDFromHex(payload.SurveyID)

	if err != nil {
		return nil, errors.ErrInternalServer
	}

	anonyPassword, err := payload.AnonymousPassword()

	if err != nil {
		return nil, err
	}

	doc := &model.Answer{
		Content:        payload.Content,
		SurveyID:       surveyObjectID,
		AuthorID:       authorID,
		IsAnonymous:    payload.IsAnonymous,
		Anonymous:      payload.AnonymousName(),
		AnswerPassword: anonyPassword,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
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

func (s *AnswerService) GetFilteredAnswers(ctx context.Context, filter bson.M, opts ...options.Lister[options.FindOptions]) ([]*model.Answer, error) {
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

	return s.answerRepo.FindByFilter(ctx, m, opts...)
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

	return s.answerRepo.RateUp(ctx, answer.ID.Hex())
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
