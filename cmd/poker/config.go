package main

import (
	"poker"

	"github.com/joho/godotenv"
)

var appConfig struct {
	Mode   string `env:"MODE" default:"server"`
	AppURL string `env:"APP_URL,required"`
	Auth0  struct {
		CallbackURL  string `env:"AUTH0_CALLBACK_URL,required"`
		ClientID     string `env:"AUTH0_CLIENT_ID,required"`
		ClientSecret string `ssm:"/poker/auth0-client-secret,required"`
		Domain       string `env:"AUTH0_DOMAIN,required"`
	}
	Session struct {
		Key string `ssm:"/poker/session-key,required"`
	}
	Environment poker.Environment `env:"ENVIRONMENT,required"`
	Server      struct {
		Port string `env:"SERVER_PORT" default:"8080"`
	}
	Audio struct {
		S3Bucket string `env:"POKER_AUDIO_CACHE_BUCKET,required"`
	}
}

func loadConfig() {

	_ = godotenv.Load()

}
