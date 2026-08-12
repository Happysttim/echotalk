package service

import (
	"context"
	"echotalk/internal/auth"
	"echotalk/internal/auth/oauth"
	"echotalk/internal/errors"
	"echotalk/internal/model"
	"echotalk/internal/utils"
)

type LoginResponse struct {
	User         *model.User
	AccessToken  string
	RefreshToken string
}

type AuthService struct {
	userService  *UserService
	googleClient *oauth.GoogleClient
}

func NewAuthService(userService *UserService) *AuthService {
	return &AuthService{
		userService:  userService,
		googleClient: oauth.NewGoogleClient(),
	}
}

func (s *AuthService) findLocalUser(ctx context.Context, email string) (*model.User, error) {
	user, err := s.userService.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) findOAuthUser(ctx context.Context, authProvider model.AuthProvider, providerID string) (*model.User, error) {
	user, err := s.userService.GetUserByProvider(ctx, authProvider, providerID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) LoginWithGoogle(ctx context.Context, payload *model.GoogleAuthRequest) (*LoginResponse, error) {
	if payload == nil {
		return nil, errors.ErrInvalidInput
	}

	oauthResult, err := s.googleClient.GoogleExchange(ctx, payload.Code)
	if err != nil {
		return nil, err
	}

	if oauthResult.Claims == nil {
		return nil, errors.ErrInvalidOAuthClaims
	}

	user, err := s.findOAuthUser(ctx, model.AuthProviderGoogle, oauthResult.Claims.Sub)

	if err != nil {
		return nil, errors.ErrInternalServer
	}

	if user == nil {
		command := &model.CreateUserCommand{
			AuthProvider: model.AuthProviderGoogle,
			ProviderID:   oauthResult.Claims.Sub,
			Email:        oauthResult.Claims.Email,
		}

		user, err = s.userService.CreateUser(ctx, command)

		if err != nil {
			return nil, err
		}

	}

	jwtToken, err := auth.CreateToken(user.ID.Hex())
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		User:         user,
		AccessToken:  jwtToken.AccessToken,
		RefreshToken: jwtToken.RefreshToken,
	}, nil
}

func (s *AuthService) RegisterWithLocal(ctx context.Context, payload *model.LocalAuthRequest) (*LoginResponse, error) {
	if payload == nil {
		return nil, errors.ErrInvalidInput
	}

	exists, err := s.userService.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	if exists != nil {
		return nil, errors.ErrBadRequest
	}

	passwordHash, err := utils.HashPassword(payload.Password)

	if err != nil {
		return nil, errors.ErrInternalServer
	}

	command := &model.CreateUserCommand{
		AuthProvider: model.AuthProviderLocal,
		Email:        payload.Email,
		PasswordHash: passwordHash,
	}

	user, err := s.userService.CreateUser(ctx, command)

	if err != nil {
		return nil, err
	}

	jwtToken, err := auth.CreateToken(user.ID.Hex())
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		User:         user,
		AccessToken:  jwtToken.AccessToken,
		RefreshToken: jwtToken.RefreshToken,
	}, nil
}

func (s *AuthService) LoginWithLocal(ctx context.Context, payload *model.LocalAuthRequest) (*LoginResponse, error) {
	if payload == nil {
		return nil, errors.ErrInvalidInput
	}

	passwordHash, err := utils.HashPassword(payload.Password)

	if err != nil {
		return nil, errors.ErrInternalServer
	}

	user, err := s.userService.GetUserByLocal(ctx, payload.Email, passwordHash)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	if user == nil {
		return nil, errors.ErrInvalidUser
	}

	jwtToken, err := auth.CreateToken(user.ID.Hex())
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		User:         user,
		AccessToken:  jwtToken.AccessToken,
		RefreshToken: jwtToken.RefreshToken,
	}, nil
}
