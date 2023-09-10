package server

import (
	"fmt"
	"math"
	"net/http"
	"poker/internal"
	"poker/internal/templates"
	"strconv"
	"strings"

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

	if len(timer.Levels) <= 0 {
		location, err := s.BuildRoute("dashboard-timer", "timerID", timer.ID)
		if err != nil {
			s.logger.WithError(err).Error("failed to build route to redirect to")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Location", location)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	var currentLevel = int(timer.CurrentLevel)
	if currentLevel >= len(timer.Levels)-1 {
		timer.CurrentLevel = 0
	}

	level := timer.Levels[currentLevel]

	level.DurationStr = formatDuration(int(level.DurationSec))

	err = s.templates.Play(ctx, &templates.PlayProps{
		User:         internal.UserFromContext(ctx),
		Timer:        timer,
		Level:        level,
		CurrentLevel: timer.CurrentLevel + 1,
	}).Render(w)
	if err != nil {
		s.logger.WithError(err).Error("failed to render dashboard timer")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (s *server) handleGetPlayTimerResetLevel(w http.ResponseWriter, r *http.Request) {

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

	timer.IsComplete = false

	err = s.timerRepo.SaveTimer(ctx, timer)
	if err != nil {
		s.logger.WithError(err).Error("failed to save timer")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	level := timer.Levels[timer.CurrentLevel]

	level.DurationStr = formatDuration(int(level.DurationSec))

	w.Header().Set("HX-Trigger-After-Settle", "countdown::reset")
	err = s.templates.TimerMasthead(
		ctx,
		timer,
		level,
	).Render(w)
	if err != nil {
		s.logger.WithError(err).Error("failed to render dashboard timer")
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func (s *server) handleGetPlayTimerNextLevel(w http.ResponseWriter, r *http.Request) {

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

	level := timer.Levels[timer.CurrentLevel]

	if timer.CurrentLevel == uint(len(timer.Levels)-1) {

		timer.IsComplete = true

		err = s.templates.TimerMasthead(
			ctx,
			timer,
			level,
		).Render(w)
		if err != nil {
			s.logger.WithError(err).Error("failed to render dashboard timer")
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	timer.CurrentLevel += 1

	err = s.timerRepo.SaveTimer(ctx, timer)
	if err != nil {
		s.logger.WithError(err).Error("failed to save timer")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	level = timer.Levels[timer.CurrentLevel]

	level.DurationStr = formatDuration(int(level.DurationSec))

	var proceed = false
	proceedStr := r.URL.Query().Get("proceed")
	if proceedStr != "" {
		parsedProceed, err := strconv.ParseBool(proceedStr)
		if err == nil {
			proceed = parsedProceed
		}
	}

	if proceed {
		w.Header().Set("HX-Trigger-After-Settle", "countdown::proceed")
	} else {
		w.Header().Set("HX-Trigger-After-Settle", "countdown::reset")
	}
	err = s.templates.TimerMasthead(
		ctx,
		timer,
		level,
	).Render(w)
	if err != nil {
		s.logger.WithError(err).Error("failed to render dashboard timer")
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func (s *server) handleGetPlayTimerPreviousLevel(w http.ResponseWriter, r *http.Request) {

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

	level := timer.Levels[timer.CurrentLevel]

	if timer.CurrentLevel == 0 {
		err = s.templates.TimerMasthead(
			ctx,
			timer,
			level,
		).Render(w)
		if err != nil {
			s.logger.WithError(err).Error("failed to render dashboard timer")
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	timer.CurrentLevel -= 1
	timer.IsComplete = false

	err = s.timerRepo.SaveTimer(ctx, timer)
	if err != nil {
		s.logger.WithError(err).Error("failed to save timer")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	level = timer.Levels[timer.CurrentLevel]

	level.DurationStr = formatDuration(int(level.DurationSec))

	w.Header().Set("HX-Trigger-After-Settle", "countdown::reset")
	err = s.templates.TimerMasthead(
		ctx,
		timer,
		level,
	).Render(w)
	if err != nil {
		s.logger.WithError(err).Error("failed to render dashboard timer")
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func formatDuration(duration int) string {

	hours := math.Floor(math.Mod(float64(duration/(60*60)), 24))
	minutes := math.Floor(math.Mod(float64(duration/60), 60))
	seconds := math.Floor(math.Mod(float64(duration), 60))

	bits := []string{}
	if hours > 0 {
		bits = append(bits, fmt.Sprintf("%02.f", hours))
	}

	bits = append(
		bits,
		fmt.Sprintf("%02.f", minutes),
		fmt.Sprintf("%02.f", seconds),
	)

	return strings.Join(bits, ":")

}
