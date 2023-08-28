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

	s.templates = tmpl

	s.logger.Infof("Starting Server on Port %s", s.port)
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

	// spew.Dump(name, pairsInf)

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
	router.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if !s.env.IsProduction() {
				s.templates.RefreshTemplates()
			}

			h.ServeHTTP(w, r)

		})
	})
	router.Use(s.user)

	router.HandleFunc("/", s.handleHome).Name("home")
	router.HandleFunc("/login", s.handleLogin).Name("login")
	router.HandleFunc("/static/style.css", s.handleCSS).Name("styles-css")
	router.PathPrefix("/static").Handler(http.StripPrefix("/static/", http.FileServer(http.FS(poker.AssetFS(s.env))))).Name("static")

	authed := router.NewRoute().Subrouter()
	authed.Use(s.auth)
	authed.HandleFunc("/dashboard", s.handleDashboard).Name("dashboard")
	authed.HandleFunc("/dashboard/timers", s.handleDashboardTimers).Name("dashboard-timers")
	authed.HandleFunc("/dashboard/timers/{timerID}", s.handleGetDashboardTimer).Methods(http.MethodGet).Name("dashboard-timer")

	partials := router.NewRoute().Subrouter()
	partials.Use(s.auth)
	partials.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			hxRequestHeader := r.Header.Get("HX-Request")
			if hxRequestHeader != "true" {
				s.logger.Error("Required Header HX-Request missing from request")
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			h.ServeHTTP(w, r)

		})
	})

	partials.HandleFunc("/partials/dashboard/timers/new", s.handleGetPartialDashboardNewTimer).Methods(http.MethodGet).Name("partials-dashboard-timers-new")
	partials.HandleFunc("/partials/dashboard/timers/new", s.handlePostPartialDashboardNewTimer).Methods(http.MethodPost)

	partials.HandleFunc("/partials/dashboard/timers", s.handlePartialDashboardTimers).
		Methods(http.MethodGet).Name("partials-dashboard-timers")

	partials.HandleFunc("/partials/dashboard/timers/{timerID}", s.handleGetPartialDashboardTimer).
		Methods(http.MethodGet).Name("partials-dashboard-timer")

	partials.HandleFunc("/partials/dashboard/timers/{timerID}", s.handleDeletePartialDashboardTimer).
		Methods(http.MethodDelete).Name("partials-dashboard-timer")

	partials.HandleFunc("/partials/dashboard/timers/{timerID}/levels/{levelType}", s.handleGetPartialDashboardTimerLevel).
		Methods(http.MethodGet).Name("partials-dashboard-timer-level")
	partials.HandleFunc("/partials/dashboard/timers/{timerID}/levels/{levelType}", s.handlePostPartialsDashboardTimerLevel).
		Methods(http.MethodPost)
	partials.HandleFunc("/partials/dashboard/timers/{timerID}/levels/{levelType}", s.handlePutPartialsDashboardTimerLevel).
		Methods(http.MethodPut)

	// Dashboard Standings
	partials.HandleFunc("/partials/dashboard/standings", s.handlePartialDashboardStandings).Name("partials-dashboard-standings")

	return router
}

func (s *server) handleCSS(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	buffer, err := s.templates.RenderStyles(ctx)
	if err != nil {
		s.logger.WithError(err).Error("failed to render styles")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/css; charset=utf-8")
	n, err := buffer.WriteTo(w)
	if err != nil {
		s.logger.WithError(err).Error("failed to write template to writer")
		w.WriteHeader(http.StatusInternalServerError)
	}

	s.logger.Debugf("wrote %d bytes", n)

}
