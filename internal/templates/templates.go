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

type breadcrumb struct {
	Text  string
	Route string
}

func (s *Service) newBreadcrumb(text, route string, args ...any) *breadcrumb {
	compiledRoute := s.buildRoute(route, args...)

	return &breadcrumb{
		Text:  text,
		Route: compiledRoute,
	}

}

const pathCrumbTmpl = `<div class="mx-2"><a href="%s" class="fw-bold text-dark-emphasis text-decoration-none">%s</a></div>`
const activeCrumbTmpl = `<div class="mx-2 active" aria-current="page">%s</div>`

func pathCrumb(crumb *breadcrumb) string {
	return fmt.Sprintf(pathCrumbTmpl, crumb.Route, crumb.Text)
}

func activeCrumb(crumb *breadcrumb) string {
	return fmt.Sprintf(activeCrumbTmpl, crumb.Text)
}

func (s *Service) breadcrumbs(crumbs ...*breadcrumb) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {

		_, _ = io.WriteString(w, `<div class="container">`)
		_, _ = io.WriteString(w, `<div class="row mx-2 my-2 bg-secondary-subtle rounded-pill">`)
		_, _ = io.WriteString(w, `<div class="col">`)
		_, _ = io.WriteString(w, `<div class="d-flex align-items-center fs-5" style="height: 40px">`)

		for i, crumb := range crumbs {
			active := i == len(crumbs)-1
			if active {
				_, err := io.WriteString(w, activeCrumb(crumb))
				if err != nil {
					return err
				}
				continue
			}

			_, _ = io.WriteString(w, pathCrumb(crumb))
			_, _ = io.WriteString(w, `<div class="mx-2">></div>`)

		}

		_, _ = io.WriteString(w, `</div>`)
		_, _ = io.WriteString(w, `</div>`)
		_, _ = io.WriteString(w, `</div>`)
		_, _ = io.WriteString(w, `</div>`)

		return nil

	})

}
