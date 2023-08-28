package templates

import (
	"bytes"
	"context"
	"fmt"
)

func (s *Service) RenderStyles(ctx context.Context) (*bytes.Buffer, error) {
	b := new(bytes.Buffer)
	err := s.getRegistry().ExecuteTemplate(b, "stylesheet", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to render stylesheet: %w", err)
	}

	return b, nil

}
