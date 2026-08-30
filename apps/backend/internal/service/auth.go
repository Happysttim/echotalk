package service

import (
	"context"
	"echotalk/internal/auth"
	"echotalk/internal/auth/oauth"
	"echotalk/internal/errors"
	"echotalk/internal/model"
	"echotalk/internal/repositories"
	"echotalk/internal/utils"
	"time"
)

type LoginResponse struct {
	User         *model.User
	AccessToken  string
	RefreshToken string
}

type AuthService struct {
	sessionRepo  repositories.SessionRepository
	userService  *UserService
	googleClient *oauth.GoogleClient
}

func NewAuthService(sessionRepo repositories.SessionRepository, userService *UserService) *AuthService {
	return &AuthService{
		sessionRepo:  sessionRepo,
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

func (s *AuthService) createSession(ctx context.Context, refreshToken string, expiresIn time.Time) error {
	hashedToken, err := utils.Hash(refreshToken)
	if err != nil {
		return err
	}

	session := model.Session{
		RefreshTokenHash: hashedToken,
		Revoked:          false,
		ExpiredAt:        expiresIn,
		CreatedAt:        time.Now(),
	}

	_, err = s.sessionRepo.Create(ctx, &session)

	if err != nil {
		return err
	}

	return nil
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

	s.createSession(ctx, jwtToken.RefreshToken, jwtToken.RefreshTokenExpiresIn)

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

	passwordHash, err := utils.Hash(payload.Password)

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

	s.createSession(ctx, jwtToken.RefreshToken, jwtToken.RefreshTokenExpiresIn)

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

	passwordHash, err := utils.Hash(payload.Password)

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

	s.createSession(ctx, jwtToken.RefreshToken, jwtToken.RefreshTokenExpiresIn)

	return &LoginResponse{
		User:         user,
		AccessToken:  jwtToken.AccessToken,
		RefreshToken: jwtToken.RefreshToken,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return errors.ErrInvalidInput
	}

	hashedToken, err := utils.Hash(refreshToken)
	if err != nil {
		return err
	}
	session, err := s.sessionRepo.FindByToken(ctx, hashedToken)

	if err != nil {
		return err
	}

	if session == nil {
		return errors.ErrInvalidToken
	}

	if session.ExpiredAt.Before(time.Now()) || session.Revoked {
		return errors.ErrExpiredToken
	}

	session.Revoked = true

	if err := s.sessionRepo.Update(ctx, session); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) IsUselessRefreshToken(ctx context.Context, refreshToken string) bool {
	if refreshToken == "" {
		return true
	}

	hashedToken, err := utils.Hash(refreshToken)
	if err != nil {
		return true
	}
	session, err := s.sessionRepo.FindByToken(ctx, hashedToken)

	if err != nil || session == nil {
		return true
	}

	return session.ExpiredAt.Before(time.Now()) || session.Revoked
}
