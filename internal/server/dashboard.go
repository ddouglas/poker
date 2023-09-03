package server

import (
	"net/http"
	"poker/internal"
)

func (s *server) handleDashboard(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	var user = internal.UserFromContext(ctx)
	err := s.templates.Dashboard(ctx, user).Render(ctx, w)
	if err != nil {
		s.logger.WithError(err).Error("failed to render dashboard")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
