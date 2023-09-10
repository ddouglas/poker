package server

import (
	"net/http"
	"poker/internal"
)

func (s *server) handleHome(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	var user = internal.UserFromContext(ctx)

	err := s.templates.Homepage(ctx, user).Render(w)
	if err != nil {
		s.logger.WithError(err).Error("failed to write template to writer")
		w.WriteHeader(http.StatusInternalServerError)
	}

}
