package templates

import (
	"context"
	"poker"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

func (s *Service) Homepage(ctx context.Context, user *poker.User) g.Node {
	return h.Doctype(
		h.HTML(
			h.Lang("en"),
			s.gtop(ctx),
			h.Body(
				s.gnavbar(ctx, user),
				h.Div(
					h.Class("banner"),
					h.Div(
						h.Class("overlay"),
						h.Div(
							h.Class("mt-5"),
							h.H1(g.Text("Red | Ventures Poker")),
							h.Hr(),
							h.P(
								h.Class("text-center"),
								g.Text("Free Food, Free Drinks, Great Time"),
							),
							h.P(
								h.Class("text-center"),
								h.Button(
									h.Class("btn btn-primary"),
									g.Text("Sign Up"),
								),
							),
						),
					),
				),
				s.gbottom(),
			),
		),
	)
}
