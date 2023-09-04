package templates

import (
	"context"
	"fmt"
	"io"
	"poker"
	"poker/internal/store/dynamo"
	"time"

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

func (s *Service) setCountdownData() templ.Component {

	var scriptURI = fmt.Sprintf("%s/js/countdown.js?v=%d", s.buildRoute("static"), time.Now().Unix())

	const javascriptData = `
		<script src='%s'></script>
	`

	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {

		s := fmt.Sprintf(
			javascriptData,
			scriptURI,
		)

		_, err := io.WriteString(w, s)
		return err

	})
}
