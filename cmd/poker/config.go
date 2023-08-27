package main

import (
	"poker"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var appConfig struct {
	Auth0 struct {
		CallbackURL  string `envconfig:"AUTH0_CALLBACK_URL" required:"true"`
		ClientID     string `envconfig:"AUTH0_CLIENT_ID" required:"true"`
		ClientSecret string `envconfig:"AUTH0_CLIENT_SECRET" required:"true"`
		Domain       string `envconfig:"AUTH0_DOMAIN" required:"true"`
	}
	Session struct {
		Key string `envconfig:"SESSION_KEY" required:"true"`
	}
	Environment poker.Environment `envconfig:"ENVIRONMENT" required:"true"`
	Server      struct {
		Port string `envconfig:"SERVER_PORT" required:"true"`
	}
}

func loadConfig() {

	err := godotenv.Load()
	if err != nil {
		logger.WithError(err).Fatal("failed to load .env")
	}

	err = envconfig.Process("", &appConfig)
	if err != nil {
		logger.WithError(err).Fatal("failed to load app configuration")
	}

	if !appConfig.Environment.Valid() {
		logger.WithField("environment", string(appConfig.Environment)).Fatal("invalid value provided environment")
	}

}
