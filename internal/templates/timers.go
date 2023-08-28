package templates

import (
	"bytes"
	"context"
	"fmt"
	"poker/internal"
)

func (s *Service) RenderDashboardTimers(props DashboardTimersProps) error {

	b := new(bytes.Buffer)
	err := s.getRegistry().ExecuteTemplate(b, "pages/dashboard/timers", props)
	if err != nil {
		return fmt.Errorf("failed to render pages/dashboard/timers: %w", err)
	}

	_, err = props.Write(b.Bytes())

	return err

}

func (s *Service) RenderPartialDashboardTimers(props DashboardTimersProps) (*bytes.Buffer, error) {

	b := new(bytes.Buffer)
	err := s.getRegistry().ExecuteTemplate(b, "partials/dashboard/timers", props)
	if err != nil {
		return nil, fmt.Errorf("failed to render partials/dashboard/timers: %w", err)
	}

	return b, nil

}

func (s *Service) RenderDashboardTimer(props DashboardTimerProps) error {

	b := new(bytes.Buffer)
	err := s.getRegistry().ExecuteTemplate(b, "pages/dashboard/timer", props)
	if err != nil {
		return fmt.Errorf("failed to render pages/dashboard/timer: %w", err)
	}

	_, err = props.Write(b.Bytes())

	return err

}

func (s *Service) RenderPartialDashboardNewTimer(ctx context.Context) (*bytes.Buffer, error) {

	data := &ViewData{
		User: internal.UserFromContext(ctx),
	}

	b := new(bytes.Buffer)
	err := s.getRegistry().ExecuteTemplate(b, "partials/dashboard/timers-new", data)
	if err != nil {
		return nil, fmt.Errorf("failed to render partials/dashboard/timers-new: %w", err)
	}

	return b, nil

}

func (s *Service) RenderPartialDashboardTimer(props DashboardTimerProps) (*bytes.Buffer, error) {

	b := new(bytes.Buffer)
	err := s.getRegistry().ExecuteTemplate(b, "partials/dashboard/timer", props)
	if err != nil {
		return nil, fmt.Errorf("failed to render partials/dashboard/timer: %w", err)
	}

	return b, nil

}

func (s *Service) RenderPartialDashboardTimerLevels(ctx context.Context, timerID string) (*bytes.Buffer, error) {

	timer, err := s.timerRepo.Timer(ctx, timerID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch timers for user: %w", err)
	}

	b := new(bytes.Buffer)
	err = s.getRegistry().ExecuteTemplate(b, "partials/timer/levels", timer)
	if err != nil {
		return nil, fmt.Errorf("failed to render partials/timer/levels: %w", err)
	}

	return b, nil

}

func (s *Service) RenderPartialDashboardTimerLevelBlindNew(props TimerLevelProps) (*bytes.Buffer, error) {

	b := new(bytes.Buffer)
	err := s.getRegistry().ExecuteTemplate(b, "partials/timer/level-blind-new", props)
	if err != nil {
		return nil, fmt.Errorf("failed to %s: %w", "partials/timer/level-blind-new", err)
	}

	return b, nil

}

func (s *Service) RenderPartialDashboardTimerLevelBlindEdit(props TimerLevelEditProps) (*bytes.Buffer, error) {

	b := new(bytes.Buffer)
	err := s.getRegistry().ExecuteTemplate(b, "partials/timer/level-blind-edit", props)
	if err != nil {
		return nil, fmt.Errorf("failed to %s: %w", "partials/timer/level-blind-edit", err)
	}

	return b, nil

}
