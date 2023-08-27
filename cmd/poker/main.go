package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"net/http"
	"os"
	"os/signal"
	"survey/internal/authenticator"
	"survey/internal/server"
	"survey/internal/templates"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/ddouglas/dynastore"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

type app struct {
	authenticator *authenticator.Service
	sessions      sessions.Store
}

var (
	logger = logrus.New()
)

func main() {
	loadConfig()

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		logger.WithError(err).Fatal("failed to load aws default config")
	}

	svc := dynamodb.NewFromConfig(cfg)

	sessionStore, _ := dynastore.New(svc, dynastore.TableName("poker-sessions-us-east-1"), dynastore.PrimaryKey("ID"))

	gob.Register(make(map[string]any))

	authSrv, err := authenticator.New(&authenticator.Config{
		ClientID:     appConfig.Auth0.ClientID,
		ClientSecret: appConfig.Auth0.ClientSecret,
		Tenant:       appConfig.Auth0.Domain,
		CallbackURL:  appConfig.Auth0.CallbackURL,
	})
	if err != nil {
		logger.WithError(err).Fatal("failed to provision authenticator service")
	}

	server := server.New(
		appConfig.Environment,
		appConfig.Server.Port,
		logger,
		authSrv,
		sessionStore,
	)

	tmplConfigs := make([]templates.ConfigFunc, 0)
	tmplConfigs = append(tmplConfigs, templates.WithFunction("route", server.BuildRoute))
	for _, t := range templateFiles {
		tmplConfigs = append(tmplConfigs, templates.WithTemplate(t))
	}

	tmpl, err := templates.New(
		appConfig.Environment,
		logger,
		tmplConfigs...,
	)
	if err != nil {
		logger.WithError(err).Fatal("failed to provision template service")
	}

	// Channel to listen for errors generated by api server
	serverErrors := make(chan error, 1)

	// Channel to listen for interrupts and to run a graceful shutdown
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)

	// Start up our server
	go func() {
		serverErrors <- server.Run(tmpl)
	}()

	// Blocking until read from channel(s)
	select {
	case err := <-serverErrors:
		logger.Fatalf("error starting server: %v", err.Error())

	case <-osSignals:
		logger.Println("starting server shutdown...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := server.GracefullyShutdown(ctx)
		if err != nil {
			logger.Fatalf("error trying to shutdown http server: %v", err.Error())
		}

	}
}

func (this *app) handleLogin(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	session, _ := this.sessions.Get(r, "login-attempt")

	query := r.URL.Query()

	state := query.Get("state")
	code := query.Get("code")

	if state != "" && code != "" {
		sessionState, ok := session.Values["state"]
		if !ok {
			logger.Error("invalid session, no state stored in session")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if sessionState.(string) != state {
			logger.Error("session state does equal query state, discarding")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		token, err := this.authenticator.Exchange(ctx, code)
		if err != nil {
			logger.WithError(err).Error("failed to exchange code for token")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		idToken, err := this.authenticator.VerifyIDToken(ctx, token)
		if err != nil {
			logger.WithError(err).Error("failed to exchange code for token")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var profile = make(map[string]any)
		err = idToken.Claims(&profile)
		if err != nil {
			logger.WithError(err).Error("failed to exchange code for token")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		session.Values["access_token"] = token.AccessToken
		session.Values["profile"] = profile

		err = session.Save(r, w)
		if err != nil {
			logger.WithError(err).Error("failed to exchange code for token")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Location", "/profile")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return

	}

	state, err := generateRandomState()
	if err != nil {
		// Create an error page and redirect to that. Use session flashing to flash an internal error message of sorts
		logger.WithError(err).Error("failed to generate state for authentication request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	session.Values["state"] = state

	err = session.Save(r, w)
	if err != nil {
		logger.WithError(err).Error("failed to save state for authentication request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", this.authenticator.AuthCodeURL(state))
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}

// func handleAddTodo(w http.ResponseWriter, r *http.Request) {
// 	data := r.FormValue("item")
// 	lock.Lock()
// 	defer lock.Unlock()
// 	this.Items = append(this.Items, data)
// 	err := files.ExecuteTemplate(w, "list-item", data)
// 	if err != nil {
// 		fmt.Println("yeet", err)
// 	}

// }