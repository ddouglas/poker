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

	"github.com/ddouglas/dynastore"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/sirupsen/logrus"
)

type server struct {
	env      poker.Environment
	port     string
	http     *http.Server
	logger   *logrus.Logger
	router   *mux.Router
	sessions *dynastore.Store
	decoder  *schema.Decoder

	// Services
	authenticator *authenticator.Service
	templates     *templates.Service

	// Repositories
	timerRepo *dynamo.TimerRepository
	userRepo  *dynamo.UserRepository
}

func New(
	env poker.Environment,
	port string,
	logger *logrus.Logger,

	authenticator *authenticator.Service,
	sessions *dynastore.Store,

	timerRepo *dynamo.TimerRepository,
	userRepo *dynamo.UserRepository,
) *server {

	s := &server{
		env:      env,
		port:     port,
		logger:   logger,
		sessions: sessions,
		decoder:  schema.NewDecoder(),

		authenticator: authenticator,

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
	router.Use(s.user)

	router.HandleFunc("/", s.handleHome).Name("home")
	router.HandleFunc("/login", s.handleLogin).Name("login")
	router.PathPrefix("/static").Handler(http.StripPrefix("/static/", http.FileServer(http.FS(poker.AssetFS(s.env))))).Name("static")

	authed := router.NewRoute().Subrouter()
	authed.Use(s.auth)
	authed.HandleFunc("/dashboard", s.handleDashboard).Name("dashboard")
	authed.HandleFunc("/dashboard/timers", s.handleDashboardTimers).Name("dashboard-timers")
	authed.HandleFunc("/dashboard/timers/new", func(w http.ResponseWriter, r *http.Request) {
		map[string]http.HandlerFunc{
			http.MethodGet:  s.handleGetDashboardTimerNew,
			http.MethodPost: s.handlePostDashboardTimerNew,
		}[r.Method](w, r)
	}).Name("dashboard-timers-new")

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
	}).Methods(http.MethodGet, http.MethodDelete).Name("play-timer")

	authed.HandleFunc("/play/{timerID}/levels/next", func(w http.ResponseWriter, r *http.Request) {
		map[string]http.HandlerFunc{
			http.MethodGet: s.handleGetPlayTimerNextLevel,
		}[r.Method](w, r)
	}).Methods(http.MethodGet, http.MethodDelete).Name("play-timer-next-level")

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

	return router
}
