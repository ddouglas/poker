package templates

import (
	"fmt"
	html "html/template"
	"io"
	"survey"

	"github.com/sirupsen/logrus"
)

type Service struct {
	logger      *logrus.Logger
	environment survey.Environment
	registry    *html.Template
	templates   map[string]*Template
	funcMap     html.FuncMap
}

type Template struct {
	Name string
	Path string
}

func New(
	env survey.Environment,
	logger *logrus.Logger,
	cfgFuncs ...ConfigFunc,
) (*Service, error) {
	s := &Service{
		environment: env,
		registry:    html.New(""),
		logger:      logger,
		funcMap:     make(html.FuncMap),
		templates:   make(map[string]*Template),
	}

	for _, f := range cfgFuncs {
		err := f(s)
		if err != nil {
			return nil, err
		}
	}

	return s, nil
}

type ConfigFunc func(s *Service) error

func WithTemplate(t *Template) ConfigFunc {
	return func(s *Service) error {
		if _, ok := s.templates[t.Name]; ok {
			return fmt.Errorf("template with name %s has already been registered", t.Name)
		}

		s.templates[t.Name] = t

		s.withTemplate(t)
		return nil
	}

}

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
	s.registry = html.New("")
	for _, t := range s.templates {
		s.withTemplate(t)
	}
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

func (s *Service) mustRead(path string) string {

	fs := survey.TemplateFS(s.environment)

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

func (s *Service) getRegistry() *html.Template {
	return s.registry.Funcs(s.funcMap)
}
