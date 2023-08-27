package server

import (
	"net/http"
)

func (s *server) handleHome(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	buffer, err := s.templates.RenderHomepage(ctx)
	if err != nil {
		s.logger.WithError(err).Error("failed to render homepage")
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
