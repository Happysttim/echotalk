package service

import (
	"context"
	"echotalk/internal/errors"
	"echotalk/internal/model"
	"echotalk/internal/repositories"
	"time"
)

type UserService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, payload *model.CreateUserCommand) (*model.User, error) {
	var user *model.User
	now := time.Now()

	user = &model.User{
		AuthProvider: payload.AuthProvider,
		Email:        payload.Email,
		ProviderID:   payload.ProviderID,
		PasswordHash: payload.PasswordHash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	result, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	if id == "" {
		return nil, errors.ErrInvalidInput
	}
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	if id == "" {
		return errors.ErrInvalidInput
	}
	err := s.userRepo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	users, err := s.userRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	if email == "" {
		return nil, errors.ErrInvalidInput
	}

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUserByProvider(ctx context.Context, authProvider model.AuthProvider, providerID string) (*model.User, error) {
	if !authProvider.IsValid() || providerID == "" {
		return nil, errors.ErrInvalidInput
	}

	user, err := s.userRepo.FindByProvider(ctx, authProvider, providerID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUserByLocal(ctx context.Context, email string) (*model.User, error) {
	if email == "" {
		return nil, errors.ErrInvalidInput
	}

	user, err := s.userRepo.FindByLocal(ctx, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}
