package templates

import (
	"bytes"
	"fmt"
)

func (s *Service) RenderDashboard(props DashboardProps) error {

	b := new(bytes.Buffer)
	err := s.getRegistry().ExecuteTemplate(b, "pages/dashboard", props)
	if err != nil {
		return fmt.Errorf("failed to render pages/dashboard: %w", err)
	}

	_, err = props.Write(b.Bytes())

	return err

}
