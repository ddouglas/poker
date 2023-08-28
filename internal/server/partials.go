package server

import (
	"net/http"
	"poker/internal/templates"
)

func (s *server) handlePartialDashboardStandings(w http.ResponseWriter, r *http.Request) {

	buffer, err := s.templates.RenderPartialDashboardStandings(templates.NewDashboardProps(r.Context(), w))
	if err != nil {
		s.logger.WithError(err).Error("failed to render dashboard standings")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	n, err := buffer.WriteTo(w)
	if err != nil {
		s.logger.WithError(err).Error("failed to write template to writer")
		w.WriteHeader(http.StatusInternalServerError)
	}

	s.logger.Debugf("wrote %d bytes", n)

}
