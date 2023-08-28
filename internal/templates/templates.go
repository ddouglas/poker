package templates

import (
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"poker"
	"poker/internal/store/dynamo"
	"strings"
	text "text/template"

	"github.com/sirupsen/logrus"
)

type Service struct {
	logger      *logrus.Logger
	environment poker.Environment
	registry    *text.Template
	templates   map[string]*Template
	funcMap     text.FuncMap

	timerRepo *dynamo.TimerRepository
}

type ViewData struct {
	User *poker.User
}

type Template struct {
	Name string
	Path string
}

type structWithIdx struct {
	Idx  int
	Data any
}

func New(
	env poker.Environment,
	logger *logrus.Logger,

	timerRepo *dynamo.TimerRepository,

	cfgFuncs ...ConfigFunc,
) (*Service, error) {
	s := &Service{
		environment: env,
		registry:    text.New(""),
		logger:      logger,
		funcMap:     make(text.FuncMap),
		templates:   make(map[string]*Template),

		timerRepo: timerRepo,
	}

	for _, f := range cfgFuncs {
		err := f(s)
		if err != nil {
			return nil, err
		}
	}

	s.registerFunctions()

	s.registerTemplates()

	return s, nil
}

type ConfigFunc func(s *Service) error

func WithFunction(name string, f any) ConfigFunc {

	return func(s *Service) error {
		if _, ok := s.funcMap[name]; ok {
			return fmt.Errorf("function %s has already been registered", name)
		}

		s.funcMap[name] = f
		return nil
	}

}

func (s *Service) RefreshTemplates() {
	s.registry = text.New("")

	s.registerTemplates()
}

func (s *Service) registerFunctions() {
	s.funcMap["structWithIdx"] = func(i int, data any) structWithIdx {
		return structWithIdx{
			Idx:  i,
			Data: data,
		}
	}

	s.funcMap["renderPartialDashboardStandings"] = s.RenderPartialDashboardStandings
	s.funcMap["renderNavbar"] = s.renderNavbar
	s.funcMap["renderPartialTop"] = s.renderPartialTop
	s.funcMap["renderPartialBottom"] = s.renderPartialBottom
	s.funcMap["renderPartialDashboardUserMenu"] = s.renderPartialDashboardUserMenu
	s.funcMap["renderPartialDashboardTimers"] = s.RenderPartialDashboardTimers
	s.funcMap["renderPartialDashboardTimer"] = s.RenderPartialDashboardTimer
	s.funcMap["renderPartialTimerLevel"] = s.renderPartialTimerLevel
}

func (s *Service) withTemplate(t *Template) {

	data := s.mustRead(t.Path)

	_, err := s.getRegistry().New(t.Name).Parse(data)
	if err != nil {
		s.logger.WithError(err).
			WithField("name", t.Name).WithField("path", t.Path).
			Fatal("failed to parse template")
	}
}

func (s *Service) registerTemplates() {
	templateFS := poker.TemplateFS(s.environment)

	err := fs.WalkDir(templateFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip Directories
		if d.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)

		name := strings.TrimSuffix(path, ext)

		s.withTemplate(&Template{
			Name: name,
			Path: path,
		})

		return nil

	})
	if err != nil {
		s.logger.WithError(err).Fatal("failed to walk template fs")
	}

}

func (s *Service) mustRead(path string) string {

	fs := poker.TemplateFS(s.environment)

	f, err := fs.Open(path)
	if err != nil {
		s.logger.WithError(err).
			WithField("path", path).
			Fatal("failed to open file on filesystem")
	}

	data, err := io.ReadAll(f)
	if err != nil {
		s.logger.WithError(err).WithField("path", path).Fatal("failed to read file")
	}

	return string(data)

}

func (s *Service) getRegistry() *text.Template {
	return s.registry.Funcs(s.funcMap)
}
