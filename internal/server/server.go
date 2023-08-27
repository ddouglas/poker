package server

import (
	"context"
	"fmt"
	"net/http"
	"survey"
	"survey/internal/authenticator"
	"survey/internal/templates"
	"time"

	"github.com/ddouglas/dynastore"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type server struct {
	env      survey.Environment
	port     string
	http     *http.Server
	logger   *logrus.Logger
	router   *mux.Router
	sessions *dynastore.Store

	// Services
	authenticator *authenticator.Service
	templates     *templates.Service
}

func New(
	env survey.Environment,
	port string,
	logger *logrus.Logger,

	authenticator *authenticator.Service,
	sessions *dynastore.Store,
) *server {
	s := &server{
		env:      env,
		port:     port,
		logger:   logger,
		sessions: sessions,

		authenticator: authenticator,
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

func (s *server) BuildRoute(name string, pairs ...string) (string, error) {

	route := s.router.GetRoute(name)
	if route == nil {
		return "", fmt.Errorf("failed to build url for %s with %d args, route not found", name, len(pairs))
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

	router.HandleFunc("/", s.handleHome).Name("homepage")
	router.HandleFunc("/login", s.handleLogin).Name("login")

	router.HandleFunc("/static/style.css", s.handleCSS).Name("styles-css")

	router.PathPrefix("/static").Handler(http.StripPrefix("/static/", http.FileServer(http.FS(survey.AssetFS(s.env))))).Name("static")
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
