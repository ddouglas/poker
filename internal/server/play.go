package server

import (
	"fmt"
	"math"
	"net/http"
	"poker/internal"
	"poker/internal/templates"
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

	level := timer.Levels[timer.CurrentLevel]

	level.DurationStr = formatDuration(int(level.DurationSec))
	err = s.templates.Play(ctx, &templates.PlayProps{
		User:         internal.UserFromContext(ctx),
		Timer:        timer,
		Level:        level,
		CurrentLevel: timer.CurrentLevel + 1,
	}).Render(ctx, w)
	if err != nil {
		s.logger.WithError(err).Error("failed to render dashboard timer")
		w.WriteHeader(http.StatusInternalServerError)
		return
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
		level.DurationStr = "All Done!!"
		err = s.templates.TimerMasthead(
			timer,
			level,
			timer.CurrentLevel+1,
		).Render(ctx, w)
		if err != nil {
			s.logger.WithError(err).Error("failed to render dashboard timer")
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	// timer.CurrentLevel += 1

	err = s.timerRepo.SaveTimer(ctx, timer)
	if err != nil {
		s.logger.WithError(err).Error("failed to save timer")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	level = timer.Levels[timer.CurrentLevel]

	level.DurationStr = formatDuration(int(level.DurationSec))
	err = s.templates.TimerMasthead(
		timer,
		level,
		timer.CurrentLevel+1,
	).Render(ctx, w)
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
