package templates

import (
	"bytes"
	"context"
	"fmt"
)

func (s *Service) RenderStyles(ctx context.Context) (*bytes.Buffer, error) {
	b := new(bytes.Buffer)
	err := s.getRegistry().ExecuteTemplate(b, "style-css", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to render homepage: %w", err)
	}

	return b, nil

}
