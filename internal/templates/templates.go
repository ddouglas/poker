package templates

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"poker"
	"poker/internal/store/dynamo"

	"github.com/a-h/templ"
	"github.com/sirupsen/logrus"
)

type Service struct {
	logger      *logrus.Logger
	environment poker.Environment

	funcs struct {
		buildRoute func(string, ...any) (string, error)
	}

	timerRepo *dynamo.TimerRepository
}

type ViewData struct {
	User *poker.User
}

type Template struct {
	Name string
	Path string
}

func New(
	env poker.Environment,
	logger *logrus.Logger,

	timerRepo *dynamo.TimerRepository,

) (*Service, error) {
	s := &Service{
		environment: env,
		logger:      logger,
		timerRepo:   timerRepo,
	}

	return s, nil
}

func (s *Service) SetRouteBuild(b func(string, ...any) (string, error)) {
	s.funcs.buildRoute = b
}

func (s *Service) buildRoute(name string, args ...any) string {
	route, err := s.funcs.buildRoute(name, args...)
	if err != nil {
		s.logger.WithField("name", name).WithError(err).Error("failed to generate styles-css route")
		return ""
	}

	return route
}

func (s *Service) setCountdownData(levels []*poker.TimerLevel) templ.Component {

	const javascriptData = `
		<script>
			const countdownServerData = JSON.parse('%v');
		</script>
	`

	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {

		data, err := json.Marshal(levels)
		if err != nil {
			return err
		}

		_, err = io.WriteString(w, fmt.Sprintf(javascriptData, string(data)))
		return err
	})
}
