package server

import (
	"net/http"
	"poker/internal/templates"
)

func (s *server) handleDashboard(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	props := templates.NewDashboardProps(ctx, w)

	err := s.templates.RenderDashboard(props)
	if err != nil {
		s.logger.WithError(err).Error("failed to render dashboard")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
