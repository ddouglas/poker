package server

import (
	"bytes"
	"net/http"
	"poker"
	"poker/internal"
	"poker/internal/templates"
	"strconv"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (s *server) handleDashboardTimers(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()

	user := internal.UserFromContext(ctx)

	timers, err := s.timerRepo.TimersByUserID(ctx, user.ID)
	if err != nil {
		s.logger.WithError(err).Error("failed to timers by user id")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	err = s.templates.RenderDashboardTimers(templates.NewDashboardTimersProps(ctx, timers, w))
	if err != nil {
		s.logger.WithError(err).Error("failed to render timers")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (s *server) handleGetDashboardTimer(w http.ResponseWriter, r *http.Request) {

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

	err = s.templates.RenderDashboardTimer(templates.NewDashboardTimerProps(ctx, timer, w))
	if err != nil {
		s.logger.WithError(err).Error("failed to render timer")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (s *server) handlePartialDashboardTimers(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	user := internal.UserFromContext(ctx)

	timers, err := s.timerRepo.TimersByUserID(ctx, user.ID)
	if err != nil {
		s.logger.WithError(err).Error("failed to timers by user id")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	uri, _ := s.router.Get("dashboard-timers").URL()
	w.Header().Set("HX-Push", uri.String())

	buffer, err := s.templates.RenderPartialDashboardTimers(templates.NewDashboardTimersProps(ctx, timers, w))
	if err != nil {
		s.logger.WithError(err).Error("failed to render timers")
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

func (s *server) handleGetPartialDashboardNewTimer(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	buffer, err := s.templates.RenderPartialDashboardNewTimer(ctx)
	if err != nil {
		s.logger.WithError(err).Error("failed to render partial dashboard timer form")
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

func (s *server) handlePostPartialDashboardNewTimer(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	err := r.ParseForm()
	if err != nil {
		s.logger.WithError(err).Error("failed to parse request form")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var timer = new(poker.Timer)
	err = s.decoder.Decode(timer, r.PostForm)
	if err != nil {
		s.logger.WithError(err).Error("failed to decode request form")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user := internal.UserFromContext(ctx)

	timer.ID = uuid.New().String()
	timer.UserID = user.ID

	err = s.timerRepo.SaveTimer(ctx, timer)
	if err != nil {
		s.logger.WithError(err).Error("failed to save timer")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	buffer, err := s.templates.RenderPartialDashboardTimer(templates.NewDashboardTimerProps(ctx, timer, nil))
	if err != nil {
		s.logger.WithError(err).Error("failed to render partial dashboard timer")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	uri, _ := s.router.Get("dashboard-timer").URL("timerID", timer.ID)
	w.Header().Set("HX-Push", uri.String())
	n, err := buffer.WriteTo(w)
	if err != nil {
		s.logger.WithError(err).Error("failed to write template to writer")
		w.WriteHeader(http.StatusInternalServerError)
	}

	s.logger.Debugf("wrote %d bytes", n)

}

func (s *server) handleGetPartialDashboardTimer(w http.ResponseWriter, r *http.Request) {

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
		s.logger.Error("failed to fetch timer by timerID")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	buffer, err := s.templates.RenderPartialDashboardTimer(templates.NewDashboardTimerProps(ctx, timer, nil))
	if err != nil {
		s.logger.WithError(err).Error("failed to render homepage")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	uri, _ := s.router.Get("dashboard-timer").URL("timerID", timerID)
	w.Header().Set("HX-Push", uri.String())
	n, err := buffer.WriteTo(w)
	if err != nil {
		s.logger.WithError(err).Error("failed to write template to writer")
		w.WriteHeader(http.StatusInternalServerError)
	}

	s.logger.Debugf("wrote %d bytes", n)

}

func (s *server) handleDeletePartialDashboardTimer(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	vars := mux.Vars(r)

	timerID, ok := vars["timerID"]
	if !ok {
		s.logger.Error("var timerID missing from request context")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := s.timerRepo.DeleteTimer(ctx, timerID)
	if err != nil {
		s.logger.WithError(err).Error("failed to delete timer")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user := internal.UserFromContext(ctx)

	timers, err := s.timerRepo.TimersByUserID(ctx, user.ID)
	if err != nil {
		s.logger.WithError(err).Error("failed to timers by user id")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	buffer, err := s.templates.RenderPartialDashboardTimers(templates.NewDashboardTimersProps(ctx, timers, nil))
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

func (s *server) handleGetPartialDashboardTimerLevel(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	vars := mux.Vars(r)

	timerID, ok := vars["timerID"]
	if !ok {
		s.logger.Error("var timerID missing from request context")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	levelType, ok := vars["levelType"]
	if !ok {
		s.logger.Error("var timerID missing from request context")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var levelIdx *int64
	levelIdxStr := r.URL.Query().Get("idx")
	if levelIdxStr != "" {
		i, err := strconv.ParseInt(levelIdxStr, 10, 32)
		if err != nil {
			s.logger.Error("var timerID missing from request context")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		levelIdx = &i
	}

	timer, err := s.timerRepo.Timer(ctx, timerID)
	if err != nil {
		s.logger.Error("failed to fetch timer")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	spew.Dump(timer, levelType, levelIdx)

	var buffer *bytes.Buffer

	if levelType == "blind" {

		if levelIdx != nil {
			buffer, err = s.templates.RenderPartialDashboardTimerLevelBlindEdit(templates.NewTimerLevelEditProps(ctx, timer, int(*levelIdx)))
			if err != nil {
				s.logger.WithError(err).Error("failed to render dashboard timer level blind edit template")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			buffer, err = s.templates.RenderPartialDashboardTimerLevelBlindNew(templates.NewTimerLevelProps(ctx, timer))
			if err != nil {
				s.logger.WithError(err).Error("failed to render dashboard timer level blind new template")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}

	if buffer == nil {
		s.logger.WithError(err).Error("failed to fill buffer with data")
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

func (s *server) handlePostPartialsDashboardTimerLevel(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	vars := mux.Vars(r)

	timerID, ok := vars["timerID"]
	if !ok {
		s.logger.Error("var timerID missing from request context")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	levelType, ok := vars["levelType"]
	if !ok {
		s.logger.Error("var timerID missing from request context")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	timer, err := s.timerRepo.Timer(ctx, timerID)
	if err != nil {
		s.logger.WithError(err).Error("failed to fetch timer")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = r.ParseForm()
	if err != nil {
		s.logger.WithError(err).Error("failed to parse request form")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	level := new(poker.TimerLevel)

	err = s.decoder.Decode(level, r.PostForm)
	if err != nil {
		s.logger.WithError(err).Error("failed to decode request form")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	level.Type = poker.TimerType(levelType)
	timer.Levels = append(timer.Levels, level)

	err = s.timerRepo.SaveTimer(ctx, timer)
	if err != nil {
		s.logger.WithError(err).Error("failed to update timer")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	buffer, err := s.templates.RenderPartialDashboardTimerLevels(ctx, timerID)
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

func (s *server) handlePutPartialsDashboardTimerLevel(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	vars := mux.Vars(r)

	timerID, ok := vars["timerID"]
	if !ok {
		s.logger.Error("var timerID missing from request context")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	levelType, ok := vars["levelType"]
	if !ok {
		s.logger.Error("var timerID missing from request context")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	timer, err := s.timerRepo.Timer(ctx, timerID)
	if err != nil {
		s.logger.WithError(err).Error("failed to fetch timer")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = r.ParseForm()
	if err != nil {
		s.logger.WithError(err).Error("failed to parse request form")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	idxStr := r.PostForm.Get("Idx")
	if idxStr == "" {
		s.logger.WithError(err).Error("missing index required to identify which level is being editted")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	idx, err := strconv.ParseInt(idxStr, 10, 32)
	if err != nil {
		s.logger.WithError(err).Error("missing index required to identify which level is being editted")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(timer.Levels) > int(idx) {
		s.logger.WithError(err).Error("unknown index")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	level := timer.Levels[idx]

	err = s.decoder.Decode(level, r.PostForm)
	if err != nil {
		s.logger.WithError(err).Error("failed to decode form")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	level.Type = poker.TimerType(levelType)

	err = s.timerRepo.SaveTimer(ctx, timer)
	if err != nil {
		s.logger.WithError(err).Error("failed to decode form")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	buffer, err := s.templates.RenderPartialDashboardTimerLevels(ctx, timerID)
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
