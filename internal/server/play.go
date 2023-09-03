package server

import (
	"net/http"
	"poker/internal"
	"poker/internal/templates"

	"github.com/gorilla/mux"
)

func (s *server) handleGetPlayTimer(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	vars := mux.Vars(r)

	timerID, ok := vars["timerID"]
	if !ok {
		s.logger.Error("var timerID missing from request context")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	timer, err := s.timerRepo.Timer(ctx, timerID)
	if err != nil {
		s.logger.WithError(err).Error("failed to fetch timer")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.templates.Play(ctx, &templates.PlayProps{
		User:  internal.UserFromContext(ctx),
		Timer: timer,
	}).Render(ctx, w)
	if err != nil {
		s.logger.WithError(err).Error("failed to render dashboard timer")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
