package server

import (
	"net/http"
	"poker/internal"
	"time"

	"github.com/sirupsen/logrus"
)

func (s *server) logging(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()
		entry := s.logger.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.String(),
		})
		handler.ServeHTTP(w, r)
		entry.WithField("duration", time.Since(start)).Info("request")

	})
}

func (s *server) user(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var ctx = r.Context()

		session, err := s.sessions.Get(r, "poker-session")
		if err != nil {
			// Create an error page and redirect to that. Use session flashing to flash an internal error message of sorts
			s.logger.WithError(err).Error("failed to load session")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		userID, ok := session.Values["userID"]
		if !ok {
			handler.ServeHTTP(w, r)
			return
		}

		user, err := s.userRepo.User(ctx, userID.(string))
		if err != nil {
			s.logger.WithError(err).Error("failed to look up user by id")
			// s.writeRedirectRouteName(w, "login")
			handler.ServeHTTP(w, r)
			return
		}

		ctx = internal.ContextWithUser(ctx, user)

		handler.ServeHTTP(w, r.WithContext(ctx))

	})
}

func (s *server) auth(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var ctx = r.Context()

		user := internal.UserFromContext(ctx)
		if user == nil {
			s.logger.Error("no user found in context, redirecting")
			s.writeRedirectRouteName(w, "login")
			return
		}

		handler.ServeHTTP(w, r)

	})
}

func (s *server) writeRedirectRouteName(w http.ResponseWriter, routeName string, routePairs ...string) {

	entry := s.logger.WithField("routeName", routeName)

	route := s.router.GetRoute(routeName)
	if route == nil {
		entry.Error("no route found for name, unable to perform redirect")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	uri, err := route.URL(routePairs...)
	if err != nil {
		entry.WithError(err).Error("failed to generate uri")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", uri.String())
	w.WriteHeader(http.StatusTemporaryRedirect)

}
