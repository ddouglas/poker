package templates

import (
	"context"

	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

func (s *Service) ErrorNotFound(ctx context.Context) g.Node {
	return Doctype(
		HTML(
			Lang("en"),
			s.gtop(ctx),
			Body(
				s.gnavbar(ctx),
				Div(
					Class("container"),
					Div(
						Class("row mt-4"),
						Div(
							Class("col"),
							Div(
								Class("alert alert-danger"),
								Strong(g.Text("Sorry, The requested resource is not available")),
							),
						),
					),
				),
			),
		),
	)
}

func (s *Service) ResourceUnavailable(ctx context.Context) g.Node {
	return Doctype(
		HTML(
			Lang("en"),
			s.gtop(ctx),
			Body(
				s.gnavbar(ctx),
				Div(
					Class("container"),
					Div(
						Class("row mt-4"),
						Div(
							Class("col"),
							Div(
								Class("alert alert-danger"),
								Strong(g.Text("That resource is not available right now, please try again later")),
							),
						),
					),
				),
			),
		),
	)
}
