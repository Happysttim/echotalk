package service

import (
	"echotalk/internal/auth"
	"echotalk/internal/errors"
	"echotalk/internal/model"
	"echotalk/internal/repositories"
	"echotalk/internal/utils"
	"time"
)

type UserService struct {
	userRepo *repositories.MongoUserRepository
}

func NewUserService(userRepo *repositories.MongoUserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) CreateUser(payload *model.CreateUserRequest) (*model.User, error) {
	if payload == nil {
		return nil, errors.ErrInvalidUser
	}

	if payload.Email == "" {
		return nil, errors.ErrInvalidInput
	}

	if payload.AuthProvider == "" {
		return nil, errors.ErrInvalidProvider
	}

	if (payload.AuthProvider == model.AuthProviderGoogle && payload.Code == "") || (payload.AuthProvider == model.AuthProviderLocal && payload.Password == "") {
		return nil, errors.ErrInvalidInput
	}

	var user *model.User

	switch payload.AuthProvider {
	case model.AuthProviderLocal:
		hashedPassword, err := utils.HashPassword(payload.Password)
		if err != nil {
			return nil, err
		}
		user = &model.User{
			AuthProvider: model.AuthProviderLocal,
			Email:        payload.Email,
			PasswordHash: hashedPassword,
		}
	case model.AuthProviderGoogle:
		tokenResponse, err := auth.GoogleToken(payload.Code)
		if err != nil {
			return nil, err
		}

		user = &model.User{
			AuthProvider: model.AuthProviderGoogle,
			ProviderID:   tokenResponse.AccessToken,
			Email:        payload.Email,
		}
	default:
		return nil, errors.ErrInvalidProvider
	}

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	result, err := s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *UserService) GetUserByID(id string) (*model.User, error) {
	if id == "" {
		return nil, errors.ErrInvalidInput
	}
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) DeleteUser(id string) error {
	if id == "" {
		return errors.ErrInvalidInput
	}
	err := s.userRepo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) GetAllUsers() ([]*model.User, error) {
	users, err := s.userRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return users, nil
}
