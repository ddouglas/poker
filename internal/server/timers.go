package server

import (
	"net/http"
	"poker"
	"poker/internal"
	"poker/internal/templates"

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

	err = s.templates.DashboardTimers(ctx, &templates.DashboardTimersProps{
		User:   internal.UserFromContext(ctx),
		Timers: timers,
	}).Render(w)
	if err != nil {
		s.logger.WithError(err).Error("failed to timers by user id")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (s *server) handleGetDashboardTimerNew(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	err := s.templates.DashboardNewTimerComponent(ctx).Render(w)
	if err != nil {
		s.logger.WithError(err).Error("failed to timers by user id")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (s *server) handlePostDashboardTimerNew(w http.ResponseWriter, r *http.Request) {

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

	uri, _ := s.router.Get("dashboard-timer").URL("timerID", timer.ID)
	w.Header().Set("HX-Push", uri.String())
	err = s.templates.DashboardTimerFragment(ctx, timer).Render(w)
	if err != nil {
		s.logger.WithError(err).Error("failed to render partial dashboard timer")
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

	err = s.templates.DashboardTimer(ctx, &templates.DashboardTimerProps{
		User:  internal.UserFromContext(ctx),
		Timer: timer,
	}).Render(w)
	if err != nil {
		s.logger.WithError(err).Error("failed to render dashboard timer")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *server) handlePostDashboardTimer(w http.ResponseWriter, r *http.Request) {

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

	err = r.ParseForm()
	if err != nil {
		s.logger.WithError(err).Error("failed to parse form")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	levels := timer.Levels

	err = s.decoder.Decode(&levels, r.PostForm)
	if err != nil {
		s.logger.WithError(err).Error("failed to decode form")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (s *server) handleDeleteDashboardTimer(w http.ResponseWriter, r *http.Request) {

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
		s.logger.WithError(err).Error("failed to fetch timers")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = s.templates.DashboardTimersFragment(ctx, timers).Render(w)
	if err != nil {
		s.logger.WithError(err).Error("failed to render dashboard timer")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (s *server) handleGetDashboardTimerLevelNew(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	vars := mux.Vars(r)

	timerID, ok := vars["timerID"]
	if !ok {
		s.logger.Error("var timerID missing from request context")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	levelTypeStr := r.URL.Query().Get("type")
	if levelTypeStr == "" {
		s.logger.Error("required query param type is missing or empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	levelType := poker.TimerType(levelTypeStr)
	if !levelType.Valid() {
		s.logger.Error("invalid timer type")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := s.templates.DashboardNewTimerLevelComponent(ctx, timerID, levelType).Render(w)
	if err != nil {
		s.logger.WithError(err).Error("failed to render dashboard timer")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (s *server) handlePostDashboardTimerLevelNew(w http.ResponseWriter, r *http.Request) {

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

	err = r.ParseForm()
	if err != nil {
		s.logger.WithError(err).Error("failed to parse form")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	level := new(poker.TimerLevel)

	err = s.decoder.Decode(level, r.PostForm)
	if err != nil {
		s.logger.WithError(err).Error("failed to decode form")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	level.ID = uuid.New().String()
	level.DurationSec = level.DurationMin * 60
	level.TimerID = timerID

	timer.Levels = append(timer.Levels, level)

	err = s.timerRepo.SaveTimer(ctx, timer)
	if err != nil {
		s.logger.WithError(err).Error("failed to save timer")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = s.templates.DashboardTimerFragment(ctx, timer).Render(w)
	if err != nil {
		s.logger.WithError(err).Error("failed to render dashboard timer")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (s *server) handleGetDashboardTimerLevelEdit(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	user := internal.UserFromContext(ctx)

	vars := mux.Vars(r)

	timerID, ok := vars["timerID"]
	if !ok {
		s.logger.Error("var timerID missing from request context")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	levelID, ok := vars["levelID"]
	if !ok {
		s.logger.Error("var levelID missing from request context")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	timer, err := s.timerRepo.Timer(ctx, timerID)
	if err != nil {
		s.logger.WithError(err).Error("failed to fetch timer")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if timer.UserID != user.ID {
		err = s.templates.DashboardTimerFragment(ctx, timer).Render(w)
		if err != nil {
			s.logger.WithError(err).Error("failed to render DashboardTimerFragment")
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	var level *poker.TimerLevel
	for _, lvl := range timer.Levels {
		if lvl.ID != levelID {
			continue
		}
		level = lvl
		break
	}

	if level == nil {
		err = s.templates.DashboardTimerFragment(ctx, timer).Render(w)
		if err != nil {
			s.logger.WithError(err).Error("failed to render DashboardTimerFragment")
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	err = s.templates.DashboardEditTimerLevelComponent(ctx, level).Render(w)
	if err != nil {
		s.logger.WithError(err).Error("failed to render DashboardEditTimerLevelComponent")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *server) handlePostDashboardTimerLevelEdit(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	user := internal.UserFromContext(ctx)

	vars := mux.Vars(r)

	timerID, ok := vars["timerID"]
	if !ok {
		s.logger.Error("var timerID missing from request context")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	levelID, ok := vars["levelID"]
	if !ok {
		s.logger.Error("var levelID missing from request context")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	timer, err := s.timerRepo.Timer(ctx, timerID)
	if err != nil {
		s.logger.Error("var levelID missing from request context")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if timer.UserID != user.ID {
		err = s.templates.DashboardTimerFragment(ctx, timer).Render(w)
		if err != nil {
			s.logger.WithError(err).Error("failed to render DashboardTimerFragment")
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	var level *poker.TimerLevel
	for _, lvl := range timer.Levels {
		if lvl.ID != levelID {
			continue
		}
		level = lvl
		break
	}
	if level == nil {
		err = s.templates.DashboardTimerFragment(ctx, timer).Render(w)
		if err != nil {
			s.logger.WithError(err).Error("failed to render DashboardTimerFragment")
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	err = r.ParseForm()
	if err != nil {
		s.logger.WithError(err).Error("failed to parse form")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = s.decoder.Decode(level, r.PostForm)
	if err != nil {
		s.logger.WithError(err).Error("failed to decode form")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	level.DurationSec = level.DurationMin * 60
	level.DurationStr = ""

	err = s.timerRepo.SaveTimer(ctx, timer)
	if err != nil {
		s.logger.WithError(err).Error("failed to save timer")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = s.templates.DashboardTimerFragment(ctx, timer).Render(w)
	if err != nil {
		s.logger.WithError(err).Error("failed to render dashboard timer")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (s *server) handleDeleteDashboardTimerLevel(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	user := internal.UserFromContext(ctx)

	vars := mux.Vars(r)

	timerID, ok := vars["timerID"]
	if !ok {
		s.logger.Error("var timerID missing from request context")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	levelID, ok := vars["levelID"]
	if !ok {
		s.logger.Error("var levelID missing from request context")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	timer, err := s.timerRepo.Timer(ctx, timerID)
	if err != nil {
		s.logger.Error("var levelID missing from request context")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if timer.UserID != user.ID {
		err = s.templates.DashboardTimer(ctx, &templates.DashboardTimerProps{
			User:  internal.UserFromContext(ctx),
			Timer: timer,
		}).Render(w)
		if err != nil {
			s.logger.WithError(err).Error("failed to render dashboard timer")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	var i int
	for idx, level := range timer.Levels {
		if level.ID != levelID {
			continue
		}
		i = idx
	}

	timer.Levels = trimLevel(timer.Levels, i)

	err = s.timerRepo.SaveTimer(ctx, timer)
	if err != nil {
		s.logger.WithError(err).Error("failed to save timer")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = s.templates.DashboardTimer(ctx, &templates.DashboardTimerProps{
		User:  internal.UserFromContext(ctx),
		Timer: timer,
	}).Render(w)
	if err != nil {
		s.logger.WithError(err).Error("failed to render dashboard timer")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func trimLevel(levels []*poker.TimerLevel, i int) []*poker.TimerLevel {

	// if i greater than the total number of levels, just bail and return the original
	if i > len(levels) {
		return levels
	}

	// if i equal the last index, return all but the last index
	if i == len(levels)-1 {
		return levels[:i]
	}

	return append(levels[:i], levels[i+1:]...)

}

// func (s *server) handlePartialDashboardTimers(w http.ResponseWriter, r *http.Request) {

// var ctx = r.Context()

// user := internal.UserFromContext(ctx)

// 	timers, err := s.timerRepo.TimersByUserID(ctx, user.ID)
// 	if err != nil {
// 		s.logger.WithError(err).Error("failed to timers by user id")
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return

// 	}

// 	uri, _ := s.router.Get("dashboard-timers").URL()
// 	w.Header().Set("HX-Push", uri.String())

// 	buffer, err := s.templates.RenderPartialDashboardTimers(templates.NewDashboardTimersProps(ctx, timers, w))
// 	if err != nil {
// 		s.logger.WithError(err).Error("failed to render timers")
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	n, err := buffer.WriteTo(w)
// 	if err != nil {
// 		s.logger.WithError(err).Error("failed to write template to writer")
// 		w.WriteHeader(http.StatusInternalServerError)
// 	}

// 	s.logger.Debugf("wrote %d bytes", n)

// }

// func (s *server) handleGetPartialDashboardTimer(w http.ResponseWriter, r *http.Request) {

// 	var ctx = r.Context()

// 	vars := mux.Vars(r)

// 	timerID, ok := vars["timerID"]
// 	if !ok {
// 		s.logger.Error("var timerID missing from request context")
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	timer, err := s.timerRepo.Timer(ctx, timerID)
// 	if err != nil {
// 		s.logger.Error("failed to fetch timer by timerID")
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	buffer, err := s.templates.RenderPartialDashboardTimer(templates.NewDashboardTimerProps(ctx, timer, nil))
// 	if err != nil {
// 		s.logger.WithError(err).Error("failed to render homepage")
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	uri, _ := s.router.Get("dashboard-timer").URL("timerID", timerID)
// 	w.Header().Set("HX-Push", uri.String())
// 	n, err := buffer.WriteTo(w)
// 	if err != nil {
// 		s.logger.WithError(err).Error("failed to write template to writer")
// 		w.WriteHeader(http.StatusInternalServerError)
// 	}

// 	s.logger.Debugf("wrote %d bytes", n)

// }

// 	user := internal.UserFromContext(ctx)

// 	timers, err := s.timerRepo.TimersByUserID(ctx, user.ID)
// 	if err != nil {
// 		s.logger.WithError(err).Error("failed to timers by user id")
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return

// 	}

// 	buffer, err := s.templates.RenderPartialDashboardTimers(templates.NewDashboardTimersProps(ctx, timers, nil))
// 	if err != nil {
// 		s.logger.WithError(err).Error("failed to render homepage")
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	n, err := buffer.WriteTo(w)
// 	if err != nil {
// 		s.logger.WithError(err).Error("failed to write template to writer")
// 		w.WriteHeader(http.StatusInternalServerError)
// 	}

// 	s.logger.Debugf("wrote %d bytes", n)

// }

// func (s *server) handleGetPartialDashboardTimerLevel(w http.ResponseWriter, r *http.Request) {

// 	var ctx = r.Context()

// 	vars := mux.Vars(r)

// 	timerID, ok := vars["timerID"]
// 	if !ok {
// 		s.logger.Error("var timerID missing from request context")
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	levelType, ok := vars["levelType"]
// 	if !ok {
// 		s.logger.Error("var timerID missing from request context")
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	var levelIdx *int64
// 	levelIdxStr := r.URL.Query().Get("idx")
// 	if levelIdxStr != "" {
// 		i, err := strconv.ParseInt(levelIdxStr, 10, 32)
// 		if err != nil {
// 			s.logger.Error("var timerID missing from request context")
// 			w.WriteHeader(http.StatusBadRequest)
// 			return
// 		}

// 		levelIdx = &i
// 	}

// 	timer, err := s.timerRepo.Timer(ctx, timerID)
// 	if err != nil {
// 		s.logger.Error("failed to fetch timer")
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	var buffer *bytes.Buffer

// 	if levelType == "blind" {

// 		if levelIdx != nil {
// 			buffer, err = s.templates.RenderPartialDashboardTimerLevelBlindEdit(templates.NewTimerLevelEditProps(ctx, timer, int(*levelIdx)))
// 			if err != nil {
// 				s.logger.WithError(err).Error("failed to render dashboard timer level blind edit template")
// 				w.WriteHeader(http.StatusInternalServerError)
// 				return
// 			}
// 		} else {
// 			buffer, err = s.templates.RenderPartialDashboardTimerLevelBlindNew(templates.NewTimerLevelProps(ctx, timer))
// 			if err != nil {
// 				s.logger.WithError(err).Error("failed to render dashboard timer level blind new template")
// 				w.WriteHeader(http.StatusInternalServerError)
// 				return
// 			}
// 		}
// 	}

// 	if buffer == nil {
// 		s.logger.WithError(err).Error("failed to fill buffer with data")
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	n, err := buffer.WriteTo(w)
// 	if err != nil {
// 		s.logger.WithError(err).Error("failed to write template to writer")
// 		w.WriteHeader(http.StatusInternalServerError)
// 	}

// 	s.logger.Debugf("wrote %d bytes", n)

// }

// func (s *server) handlePostPartialsDashboardTimerLevel(w http.ResponseWriter, r *http.Request) {

// 	var ctx = r.Context()

// 	vars := mux.Vars(r)

// 	timerID, ok := vars["timerID"]
// 	if !ok {
// 		s.logger.Error("var timerID missing from request context")
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	levelType, ok := vars["levelType"]
// 	if !ok {
// 		s.logger.Error("var timerID missing from request context")
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	timer, err := s.timerRepo.Timer(ctx, timerID)
// 	if err != nil {
// 		s.logger.WithError(err).Error("failed to fetch timer")
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	err = r.ParseForm()
// 	if err != nil {
// 		s.logger.WithError(err).Error("failed to parse request form")
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	fmt.Println(r.PostForm)

// 	level := new(poker.TimerLevel)

// 	err = s.decoder.Decode(level, r.PostForm)
// 	if err != nil {
// 		s.logger.WithError(err).Error("failed to decode request form")
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	level.TimerID = timerID

// 	level.Type = poker.TimerType(levelType)
// 	timer.Levels = append(timer.Levels, level)

// 	err = s.timerRepo.SaveTimer(ctx, timer)
// 	if err != nil {
// 		s.logger.WithError(err).Error("failed to update timer")
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	buffer, err := s.templates.RenderPartialDashboardTimerLevels(ctx, timerID)
// 	if err != nil {
// 		s.logger.WithError(err).Error("failed to render homepage")
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	n, err := buffer.WriteTo(w)
// 	if err != nil {
// 		s.logger.WithError(err).Error("failed to write template to writer")
// 		w.WriteHeader(http.StatusInternalServerError)
// 	}

// 	s.logger.Debugf("wrote %d bytes", n)

// }

// func (s *server) handlePutPartialsDashboardTimerLevel(w http.ResponseWriter, r *http.Request) {

// 	var ctx = r.Context()

// 	vars := mux.Vars(r)

// 	timerID, ok := vars["timerID"]
// 	if !ok {
// 		s.logger.Error("var timerID missing from request context")
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	levelType, ok := vars["levelType"]
// 	if !ok {
// 		s.logger.Error("var timerID missing from request context")
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	timer, err := s.timerRepo.Timer(ctx, timerID)
// 	if err != nil {
// 		s.logger.WithError(err).Error("failed to fetch timer")
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	err = r.ParseForm()
// 	if err != nil {
// 		s.logger.WithError(err).Error("failed to parse request form")
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	idxStr := r.PostForm.Get("Idx")
// 	if idxStr == "" {
// 		s.logger.WithError(err).Error("missing index required to identify which level is being editted")
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	idx, err := strconv.ParseInt(idxStr, 10, 32)
// 	if err != nil {
// 		s.logger.WithError(err).Error("missing index required to identify which level is being editted")
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	if len(timer.Levels) > int(idx) {
// 		s.logger.WithError(err).Error("unknown index")
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	level := timer.Levels[idx]

// 	err = s.decoder.Decode(level, r.PostForm)
// 	if err != nil {
// 		s.logger.WithError(err).Error("failed to decode form")
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	level.Type = poker.TimerType(levelType)

// 	err = s.timerRepo.SaveTimer(ctx, timer)
// 	if err != nil {
// 		s.logger.WithError(err).Error("failed to decode form")
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	buffer, err := s.templates.RenderPartialDashboardTimerLevels(ctx, timerID)
// 	if err != nil {
// 		s.logger.WithError(err).Error("failed to render homepage")
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	n, err := buffer.WriteTo(w)
// 	if err != nil {
// 		s.logger.WithError(err).Error("failed to write template to writer")
// 		w.WriteHeader(http.StatusInternalServerError)
// 	}

// 	s.logger.Debugf("wrote %d bytes", n)

// }
