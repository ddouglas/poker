package server

import (
	"net/http"
	"poker/internal"
)

func (s *server) auth(handler http.Handler) http.Handler {
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
			s.writeRedirectRouteName(w, "login")
			return
		}

		user, err := s.userRepo.User(ctx, userID.(string))
		if err != nil {
			s.logger.WithError(err).Error("failed to look up user by id")
			s.writeRedirectRouteName(w, "login")
			return
		}

		ctx = internal.ContextWithUser(ctx, user)

		handler.ServeHTTP(w, r.WithContext(ctx))

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
