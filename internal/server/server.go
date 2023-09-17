package server

import (
	"context"
	"fmt"
	"net/http"
	"poker"
	"poker/internal/authenticator"
	"poker/internal/store/dynamo"
	"poker/internal/templates"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/polly"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ddouglas/dynastore"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/sirupsen/logrus"
)

type server struct {
	appURL      string
	env         poker.Environment
	port        string
	audioBucket string

	http   *http.Server
	router *mux.Router

	// Services
	authenticator *authenticator.Service
	decoder       *schema.Decoder
	logger        *logrus.Logger
	polly         *polly.Client
	s3            *s3.Client
	sessions      *dynastore.Store
	templates     *templates.Service
	validator     *validator.Validate

	// Repositories
	timerRepo *dynamo.TimerRepository
	userRepo  *dynamo.UserRepository
}

func New(
	env poker.Environment,
	appURL string,
	port string,
	audioBucket string,
	logger *logrus.Logger,
	validator *validator.Validate,

	authenticator *authenticator.Service,
	polly *polly.Client,
	s3 *s3.Client,
	sessions *dynastore.Store,

	timerRepo *dynamo.TimerRepository,
	userRepo *dynamo.UserRepository,
) *server {

	s := &server{
		appURL:      appURL,
		env:         env,
		port:        port,
		audioBucket: audioBucket,

		authenticator: authenticator,
		decoder:       schema.NewDecoder(),
		logger:        logger,
		polly:         polly,
		s3:            s3,
		sessions:      sessions,
		validator:     validator,

		timerRepo: timerRepo,
		userRepo:  userRepo,
	}

	s.router = s.buildRouter()

	s.http = &http.Server{
		Addr:         fmt.Sprintf(":%s", s.port),
		ReadTimeout:  time.Second * 2,
		WriteTimeout: time.Second * 3,
		Handler:      s.router,
	}

	return s
}

func (s *server) Mux(tmpl *templates.Service) *mux.Router {
	tmpl.SetRouteBuild(s.BuildRoute)
	s.templates = tmpl

	return s.router
}

func (s *server) Run(tmpl *templates.Service) error {

	tmpl.SetRouteBuild(s.BuildRoute)
	s.templates = tmpl

	s.logger.Infof("Starting Server: http://localhost:%s", s.port)
	return s.http.ListenAndServe()
}

func (s *server) GracefullyShutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}

func (s *server) BuildRoute(name string, pairsInf ...any) (string, error) {

	route := s.router.GetRoute(name)
	if route == nil {
		return "", fmt.Errorf("failed to build url for %s with %d args, route not found", name, len(pairsInf))
	}

	var pairs = make([]string, 0, len(pairsInf))
	for _, o := range pairsInf {
		pairs = append(pairs, fmt.Sprintf("%v", o))
	}

	out, err := route.URL(pairs...)
	if err != nil {
		return "", fmt.Errorf("failed to build url for %s with %d args: %w", name, len(pairs), err)
	}

	return out.String(), nil

}

func (s *server) buildRouter() *mux.Router {

	router := mux.NewRouter()
	router.Use(s.logging)
	router.Use(s.user)

	router.HandleFunc("/", s.handleHome).Name("home").Methods(http.MethodGet)
	router.HandleFunc("/login", s.handleLogin).Name("login").Methods(http.MethodGet)
	router.HandleFunc("/logout", s.handleLogout).Name("logout").Methods(http.MethodGet)
	// router.PathPrefix("/static").Handler().Name("static").Methods(http.MethodGet)
	router.PathPrefix("/static").Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("cache-control", "max-age=86400")
		http.StripPrefix("/static/", http.FileServer(http.FS(poker.AssetFS(s.env)))).ServeHTTP(w, r)
	})).Name("static").Methods(http.MethodGet)

	authed := router.NewRoute().Subrouter()
	authed.Use(s.auth)
	authed.HandleFunc("/dashboard", s.handleDashboard).Name("dashboard").Methods(http.MethodGet)
	authed.HandleFunc("/dashboard/timers", s.handleDashboardTimers).Name("dashboard-timers").Methods(http.MethodGet)
	authed.HandleFunc("/dashboard/timers/new", func(w http.ResponseWriter, r *http.Request) {
		map[string]http.HandlerFunc{
			http.MethodGet:  s.handleGetDashboardTimerNew,
			http.MethodPost: s.handlePostDashboardTimerNew,
		}[r.Method](w, r)
	}).Methods(http.MethodGet, http.MethodPost).Name("dashboard-timers-new")

	authed.HandleFunc("/dashboard/timers/{timerID}", func(w http.ResponseWriter, r *http.Request) {
		map[string]http.HandlerFunc{
			http.MethodGet:    s.handleGetDashboardTimer,
			http.MethodDelete: s.handleDeleteDashboardTimer,
		}[r.Method](w, r)
	}).Methods(http.MethodGet, http.MethodDelete).Name("dashboard-timer")

	authed.HandleFunc("/play/{timerID}", func(w http.ResponseWriter, r *http.Request) {
		map[string]http.HandlerFunc{
			http.MethodGet: s.handleGetPlayTimer,
		}[r.Method](w, r)
	}).Methods(http.MethodGet).Name("play-timer")

	authed.HandleFunc("/play/{timerID}/levels/reset", func(w http.ResponseWriter, r *http.Request) {
		map[string]http.HandlerFunc{
			http.MethodGet: s.handleGetPlayTimerResetLevel,
		}[r.Method](w, r)
	}).Methods(http.MethodGet).Name("play-timer-reset-level")

	authed.HandleFunc("/play/{timerID}/levels/next", func(w http.ResponseWriter, r *http.Request) {
		map[string]http.HandlerFunc{
			http.MethodGet: s.handleGetPlayTimerNextLevel,
		}[r.Method](w, r)
	}).Methods(http.MethodGet).Name("play-timer-next-level")

	authed.HandleFunc("/play/{timerID}/levels/previous", func(w http.ResponseWriter, r *http.Request) {
		map[string]http.HandlerFunc{
			http.MethodGet: s.handleGetPlayTimerPreviousLevel,
		}[r.Method](w, r)
	}).Methods(http.MethodGet).Name("play-timer-previous-level")

	authed.HandleFunc("/dashboard/timers/{timerID}/levels/new", func(w http.ResponseWriter, r *http.Request) {
		map[string]http.HandlerFunc{
			http.MethodGet:  s.handleGetDashboardTimerLevelNew,
			http.MethodPost: s.handlePostDashboardTimerLevelNew,
		}[r.Method](w, r)
	}).Methods(http.MethodGet, http.MethodPost).Name("dashboard-timer-levels")

	authed.HandleFunc("/dashboard/timers/{timerID}/levels/{levelID}", func(w http.ResponseWriter, r *http.Request) {
		map[string]http.HandlerFunc{
			http.MethodGet:    s.handleGetDashboardTimerLevelEdit,
			http.MethodPost:   s.handlePostDashboardTimerLevelEdit,
			http.MethodDelete: s.handleDeleteDashboardTimerLevel,
		}[r.Method](w, r)
	}).Methods(http.MethodGet, http.MethodPost, http.MethodDelete).Name("dashboard-timer-level")

	authed.HandleFunc("/dashboard/timers/{timerID}/levels/{levelID}/audio/{action}", func(w http.ResponseWriter, r *http.Request) {
		map[string]http.HandlerFunc{
			http.MethodGet: s.handleGetDashboardTimerLevelAudio,
		}[r.Method](w, r)
	}).Methods(http.MethodGet, http.MethodPost, http.MethodDelete).Name("dashboard-timer-level-audio")

	return router
}
