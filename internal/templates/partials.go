package templates

import (
	"bytes"
	"fmt"
	"poker"
)

func (s *Service) RenderPartialDashboardStandings(props ContextProps) (*bytes.Buffer, error) {

	b := new(bytes.Buffer)
	err := s.getRegistry().ExecuteTemplate(b, "partials/dashboard/standings", props)
	if err != nil {
		return nil, fmt.Errorf("failed to render partials/dashboard/standings: %w", err)
	}

	return b, nil

}

func (s *Service) renderPartialTop(props ContextProps) (*bytes.Buffer, error) {

	b := new(bytes.Buffer)
	err := s.getRegistry().ExecuteTemplate(b, "partials/top", props)
	if err != nil {
		return nil, fmt.Errorf("failed to render partials/top: %w", err)
	}

	return b, nil

}

func (s *Service) renderNavbar(props NavbarProps) (*bytes.Buffer, error) {

	b := new(bytes.Buffer)
	err := s.getRegistry().ExecuteTemplate(b, "partials/navbar", props)
	if err != nil {
		return nil, fmt.Errorf("failed to render partials/navbar: %w", err)
	}

	return b, nil

}

func (s *Service) renderPartialBottom(props ContextProps) (*bytes.Buffer, error) {
	b := new(bytes.Buffer)
	err := s.getRegistry().ExecuteTemplate(b, "partials/bottom", props)
	if err != nil {
		return nil, fmt.Errorf("failed to render partials/bottom: %w", err)
	}

	return b, nil
}

func (s *Service) renderPartialDashboardUserMenu(props ContextProps) (*bytes.Buffer, error) {
	b := new(bytes.Buffer)
	err := s.getRegistry().ExecuteTemplate(b, "partials/dashboard/user-menu", props)
	if err != nil {
		return nil, fmt.Errorf("failed to render partials/dashboard/user-menu: %w", err)
	}

	return b, nil
}

func (s *Service) renderPartialTimerLevel(idx int, level *poker.TimerLevel) (*bytes.Buffer, error) {

	b := new(bytes.Buffer)

	data := struct {
		Idx  int
		Data *poker.TimerLevel
	}{
		Idx:  idx,
		Data: level,
	}

	err := s.getRegistry().ExecuteTemplate(b, "partials/timer/level", data)
	if err != nil {
		return nil, fmt.Errorf("failed to render partials/timer/level: %w", err)
	}

	return b, nil

}
