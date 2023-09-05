package authenticator

import (
	"context"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type Service struct {
	*oidc.Provider
	*oauth2.Config

	IssuerURL string
}

type Config struct {
	Tenant       string
	ClientID     string
	ClientSecret string
	CallbackURL  string
}

func New(cfg *Config) (*Service, error) {

	var issuerURL = fmt.Sprintf("https://%s/", cfg.Tenant)

	provider, err := oidc.NewProvider(
		context.Background(),
		issuerURL,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to provision oidc provider: %w", err)
	}

	return &Service{
		Provider: provider,
		Config: &oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURL:  cfg.CallbackURL,
			Endpoint:     provider.Endpoint(),
			Scopes:       []string{oidc.ScopeOpenID, "profile"},
		},
		IssuerURL: issuerURL,
	}, nil

}

func (a *Service) VerifyIDToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error) {

	raw, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, fmt.Errorf("field id_token missing from oauth2 token")
	}

	oidcConfig := &oidc.Config{
		ClientID: a.ClientID,
	}

	return a.Verifier(oidcConfig).Verify(ctx, raw)

}
