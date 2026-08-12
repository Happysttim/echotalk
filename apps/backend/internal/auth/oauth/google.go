package oauth

import (
	"context"
	"echotalk/internal/config"
	"echotalk/internal/errors"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleClaims struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

type GoogleAuthResult struct {
	Claims *GoogleClaims
	Token  *oauth2.Token
}

type GoogleClient struct {
	config *oauth2.Config
}

func NewGoogleClient() *GoogleClient {
	config := config.Config

	return &GoogleClient{
		config: &oauth2.Config{
			ClientID:     config.GoogleClientID,
			ClientSecret: config.GoogleClientSecret,
			RedirectURL:  config.GoogleRedirectURL,
			Scopes: []string{
				"openid",
				"email",
				"profile",
			},
			Endpoint: google.Endpoint,
		},
	}
}

func (g *GoogleClient) GoogleExchange(
	ctx context.Context,
	code string,
) (*GoogleAuthResult, error) {
	provider, err := oidc.NewProvider(
		ctx,
		AccountGoogle,
	)

	if err != nil {
		return nil, err
	}

	verifier := provider.Verifier(&oidc.Config{
		ClientID: config.Config.GoogleClientID,
	})

	var result GoogleAuthResult
	token, err := g.config.Exchange(ctx, code)

	if err != nil {
		return nil, err
	}

	rawToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.ErrInvalidIDToken
	}

	idToken, err := verifier.Verify(ctx, rawToken)

	if err != nil {
		return nil, err
	}

	if err := idToken.Claims(&result.Claims); err != nil {
		return nil, err
	}

	result.Token = token
	return &result, nil
}
