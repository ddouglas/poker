package templates

import (
	"bytes"
	"context"
	"fmt"
)

func (s *Service) RenderHomepage(ctx context.Context) (*bytes.Buffer, error) {

	data := HomepageProps{
		ctx: ctx,
	}

	b := new(bytes.Buffer)
	err := s.getRegistry().ExecuteTemplate(b, "pages/home", data)
	if err != nil {
		return nil, fmt.Errorf("failed to render pages/home: %w", err)
	}

	return b, nil

}
