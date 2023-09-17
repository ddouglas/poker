package templates

import (
	"context"
	"fmt"
	"poker"

	g "github.com/maragudk/gomponents"
	htmx "github.com/maragudk/gomponents-htmx"
	. "github.com/maragudk/gomponents/html"
)

func (s *Service) Dashboard(ctx context.Context, user *poker.User) g.Node {
	return Doctype(
		HTML(
			Lang("en"),
			s.gtop(ctx),
			Body(
				s.gnavbar(ctx),
				Div(
					Class("container"),
					s.dashboardUserCallout(ctx, user),
					Div(
						Class("row"),
						Div(
							Class("col-3"),
							s.dashboardUserMenuComponent(ctx),
						),
						Div(
							Class("col-9"),
							s.dashboardStandingsComponents(ctx),
						),
					),
				),
				s.gbottom(),
			),
		),
	)
}

func (s *Service) dashboardUserMenuComponent(ctx context.Context) g.Node {
	return g.Group([]g.Node{
		H5(g.Text("User Menu")),
		Hr(),
		Div(
			Class("list-group"),
			A(Href(s.buildRoute("dashboard")), Class("list-group-item list-group-item-action"), g.Text("Dashboard")),
			A(Href(s.buildRoute("dashboard-timers")), Class("list-group-item list-group-item-action"), g.Text("My Timers")),
		),
	})
}

func (s *Service) dashboardStandingsComponents(ctx context.Context) g.Node {
	return Div(
		Class("Container"), ID("dashboard-section"), g.Attr("hx-swap-oob", "true"),
		Div(
			Class("row"),
			Div(
				Class("col"),
				H5(Class("text-center"), g.Text("Your Standings")),
				Hr(),
				H2(Class("text-center"), g.Text("Coming Soon")),
			),
		),
	)
}

func (s *Service) dashboardUserCallout(ctx context.Context, user *poker.User) g.Node {
	return Div(
		Class("row"),
		Div(
			Class("col"),
			H1(
				g.Textf("Welcome %s", user.Name),
			),
			Hr(),
		),
	)
}

type DashboardTimersProps struct {
	User   *poker.User
	Timers []*poker.Timer
}

func (s *Service) DashboardTimers(ctx context.Context, props *DashboardTimersProps) g.Node {
	return Doctype(
		HTML(
			Lang("en"),
			s.gtop(ctx),
			Body(
				s.gnavbar(ctx),
				Div(
					Class("container"),
					s.dashboardUserCallout(ctx, props.User),
					Div(
						Class("row"),
						Div(
							Class("col-3"),
							s.dashboardUserMenuComponent(ctx),
						),
						Div(
							Class("col-9"),
							s.DashboardTimersFragment(ctx, props.Timers),
						),
					),
				),

				s.gbottom(),
			),
		),
	)

}

func (s *Service) DashboardTimersFragment(ctx context.Context, timers []*poker.Timer) g.Node {
	return s.dashboardTimersFragment(ctx, timers)
}

func (s *Service) dashboardTimersFragment(ctx context.Context, timers []*poker.Timer) g.Node {
	return Div(
		ID("dashboard-section"), g.Attr("hx-swap-oob", "true"),
		Div(
			Class("row"),
			Div(
				Class("col"),
				H5(Class("text-center"), g.Text("My Blind Timers")),
				Hr(),
			),
			Div(
				Class("row mb-3"),
				Div(
					Class("col"),
					Div(
						Class("list-group"),
						g.If(
							len(timers) > 0,
							func() g.Node {
								group := make([]g.Node, 0, len(timers))
								for _, timer := range timers {
									group = append(group, s.dashboardTimerListItemFragment(ctx, timer))
								}

								return g.Group(group)
							}(),
						),
						g.If(
							len(timers) == 0,
							Div(
								Class("alert alert-info text-center"),
								g.Text("You don't have any timers. Click below to create one now"),
							),
						),
					),
					Div(
						Class("row"),
						Div(
							Class("col"),
							Div(
								Class("d-flex justify-content-center mt-2"),
								Button(
									Class("btn btn-primary"), htmx.Get(s.buildRoute("dashboard-timers-new")), htmx.Target("#dashboard-section"),
									g.Text("Create New Timer"),
								),
							),
						),
					),
				),
			),
		),
	)
}

func (s *Service) dashboardTimerListItemFragment(ctx context.Context, timer *poker.Timer) g.Node {

	return Div(
		Class("list-group-item"),
		Div(
			Class("d-flex justify-content-between"),
			Div(g.Text(timer.Name)),
			Div(
				Div(
					Class("btn-group"), Role("group"),
					A(
						Class("btn btn-sm btn-success"), Href(s.buildRoute("play-timer", "timerID", timer.ID)),
						I(Class("fa-solid fa-play")),
					),
					A(
						Class("btn btn-sm btn-info"), Href(s.buildRoute("dashboard-timer", "timerID", timer.ID)),
						I(Class("fa-solid fa-pencil")),
					),
					Button(
						Class("btn btn-sm btn-danger"), Type("button"), g.Attr("hx-delete", s.buildRoute("dashboard-timer", "timerID", timer.ID)),
						g.Attr("hx-confirm", "Are you sure you want to delete this timer?"),
						I(Class("fa-solid fa-trash")),
					),
				),
			),
		),
	)

}

type DashboardTimerNewProps struct {
	Errors []string
}

func (s *Service) DashboardNewTimerComponent(ctx context.Context, props *DashboardTimerNewProps) g.Node {
	return Div(
		ID("dashboard-section"), g.Attr("hx-swap-oob", "true"),
		Div(
			Class("row"),
			Div(
				Class("col"),
				H5(Class("text-center"), g.Text("Create New Timer")),
				Hr(),
			),
			Div(
				Class("row mb-3"),
				Div(
					Class("col-6 offset-3"),
					Div(
						Class("card"),
						Div(
							Class("card-body"),
							s.renderErrorAlert(props.Errors),
							FormEl(
								g.Attr("hx-post", s.buildRoute("dashboard-timers-new")), g.Attr("hx-target", "#dashboard-section"),
								Div(
									Class("mb-3"),
									Label(
										Class("form-label"),
										g.Text("Timer Name"),
									),
									Input(
										ID("timer-name"), htmx.Preserve("true"), Type("text"), Class("form-control"), AutoComplete("off"), Name("name"),
									),
								),
								Div(
									Class("d-flex justify-content-center"),
									Button(
										Type("submit"), Class("btn btn-primary"), g.Text("Create Timer"),
									),
								),
							),
						),
					),
				),
			),
		),
	)
}

func (s *Service) renderErrorAlert(errors []string) g.Node {

	if len(errors) == 0 {
		return nil
	}

	errorLines := make([]g.Node, 0, len(errors))
	for _, err := range errors {
		errorLines = append(errorLines, Li(g.Text(err)))
	}

	return Div(
		Class("row"),
		Div(
			Class("col"),
			Div(
				Class("alert alert-danger"),
				Strong(g.Text("The following errors were encountered whilst processing your request")),
				Ul(errorLines...),
			),
		),
	)
}

type DashboardTimerProps struct {
	User  *poker.User
	Timer *poker.Timer
}

func (s *Service) DashboardTimer(ctx context.Context, props *DashboardTimerProps) g.Node {
	return Doctype(
		HTML(
			Lang("en"),
			s.gtop(ctx),
			Body(
				s.gnavbar(ctx),
				Div(
					Class("container"),
					s.dashboardUserCallout(ctx, props.User),
					Div(
						Class("row"),
						Div(
							Class("col-3"),
							s.dashboardUserMenuComponent(ctx),
						),
						Div(
							Class("col-9"),
							s.DashboardTimerFragment(ctx, props.Timer),
						),
					),
				),
				s.gbottom(),
			),
		),
	)
}

func (s *Service) DashboardTimerFragment(ctx context.Context, timer *poker.Timer) g.Node {
	return s.dashboardTimer(ctx, timer)
}

func (s *Service) dashboardTimer(ctx context.Context, timer *poker.Timer) g.Node {

	levelNodes := make([]g.Node, 0, len(timer.Levels))
	for idx, level := range timer.Levels {
		levelNodes = append(levelNodes, s.dashboardTimerLevelComponent(ctx, idx+1, level))
	}

	return Div(
		ID("dashboard-section"), g.Attr("hx-swap-oob", "true"),
		Div(
			Class("row"),
			Div(
				Class("col"),
				H5(Class("text-center"), g.Text(timer.Name)),
				Hr(),
			),
		),
		Div(
			Class("row mb-2"),
			Div(
				Class("col"),
				Table(

					ID("levels-table"),
					Class("table table-bordered"),
					THead(
						Class("table-secondary"),
						Tr(
							Th(
								g.Text("#")),
							Th(
								Width("20%"),
								Class("text-center"),
								g.Text("Small Blind")),
							Th(
								Width("20%"),
								Class("text-center"),
								g.Text("Big Blind")),
							Th(
								Width("20%"),
								Class("text-center"),
								g.Text("Ante")),
							Th(
								Width("20%"),
								Class("text-center"),
								g.Text("Duration (minutes)"),
							),
							Th(),
						),
					),
					TBody(levelNodes...),
				),
			),
		),
		Div(
			ID("modify-container"),
			Class("row"),
			Div(
				Class("col-6 offset-3"),
				Div(
					Class("d-flex justify-content-around"),
					Button(
						Class("btn btn-primary btn-sm"),
						htmx.Get(fmt.Sprintf("%s?type=%s", s.buildRoute("dashboard-timer-levels", "timerID", timer.ID), "blind")),
						htmx.Target("#modify-container"),
						htmx.Swap("outerHTML"),
						g.Text("Add Blind"),
					),
					Button(
						Class("btn btn-primary btn-sm"),
						htmx.Get(fmt.Sprintf("%s?type=%s", s.buildRoute("dashboard-timer-levels", "timerID", timer.ID), "break")),
						htmx.Target("#modify-container"),
						g.Text("Add Break"),
					),
				),
			),
		),
		Div(
			Class("row mt-3"),
			Div(
				Class("col-6 offset-3"),
				Div(
					Class("d-flex justify-content-around"),
					A(Href(s.buildRoute("play-timer", "timerID", timer.ID)), Class("btn btn-sm btn-success"), g.Text("Start Timer")),
				),
			),
		),
	)
}

func (s *Service) dashboardTimerLevelComponent(ctx context.Context, idx int, level *poker.TimerLevel) g.Node {

	return Tr(
		g.If(
			level.Type == "blind",
			group(
				Td(g.Textf("%v", idx)),
				Td(g.Textf("%v", level.SmallBlind)),
				Td(g.Textf("%v", level.BigBlind)),
				Td(g.Textf("%v", level.Ante)),
				Td(g.Textf("%v", level.DurationMin)),
			),
		),
		g.If(
			level.Type == "break",
			group(
				Td(g.Textf("%v", idx)),
				Td(
					ColSpan("3"), Class("text-center"),
					Strong(Em(g.Text("BREAK!"))),
				),
				Td(g.Textf("%v", level.DurationMin)),
			),
		),
		Td(
			Button(
				htmx.Get(s.buildRoute("dashboard-timer-level", "timerID", level.TimerID, "levelID", level.ID)),
				htmx.Target("#modify-container"),
				Class("btn btn-primary me-2"),
				I(Class("fa-solid fa-pencil")),
			),
			Button(
				htmx.Delete(s.buildRoute("dashboard-timer-level", "timerID", level.TimerID, "levelID", level.ID)),
				Class("btn btn-danger"),
				I(Class("fa-solid fa-trash")),
			),
		),
	)

}

type DashboardNewTimerLevelProps struct {
	TimerID   string
	LevelType poker.LevelType
	Errors    []string
}

func (s *Service) DashboardNewTimerLevelComponent(ctx context.Context, props *DashboardNewTimerLevelProps) g.Node {

	timerID, levelType, errors := props.TimerID, props.LevelType, props.Errors

	return Div(
		Class("row"),
		Div(
			Class("col"),
			Div(
				Class("card"),
				ID("card-new-level"),
				Div(
					Class("card-header text-center text-capitalize"),
					g.Textf("Create New %s", levelType.String()),
				),
				Div(
					Class("card-body"),
					s.renderErrorAlert(errors),
					FormEl(
						htmx.Post(s.buildRoute("dashboard-timer-levels", "timerID", timerID)),
						htmx.Target("#card-new-level"),
						Div(
							g.If(
								levelType == poker.LevelTypeBlind,
								group(
									Class("row row-cols-lg-3 align-items-center"),

									Div(
										Class("col-12"),
										Div(
											Label(g.Text("Small Blind")),
											Input(
												Class("form-control"), Type("number"), Name("SmallBlind"),
											),
										),
									),
									Div(
										Class("col-12"),
										Div(
											Label(g.Text("Big Blind")),
											Input(
												Class("form-control"), Type("number"), Name("BigBlind"),
											),
										),
									),
									// Div(
									// 	Class("col-12"),
									// 	Div(
									// 		Label(g.Text("Small Blind")),
									// 		Input(
									// 			Class("form-control"), Type("number"), Name("Ante"), Placeholder("Ante"),
									// 		),
									// 	),
									// ),
									Div(
										Class("col-12"),
										Div(
											Label(g.Text("Duration (minutes)")),
											Input(
												Class("form-control"), Type("number"), Name("DurationMin"),
											),
										),
									),
								),
							),
							g.If(
								levelType == poker.LevelTypeBreak,
								group(
									Class("row row-cols-lg-4 justify-content-center align-items-center"),

									Div(
										Class("col-12"),
										Div(
											Label(g.Text("Duration (minutes)")),
											Input(
												Class("form-control"), Type("number"), Name("DurationMin"),
											),
										),
									),
								),
							),
						),
						Div(
							Class("row"),
							Div(
								Class("col d-flex justify-content-center"),
								Input(
									Type("hidden"), Name("Type"), Value(levelType.String()),
								),
								Button(
									Type("submit"),
									Class("btn btn-sm btn-primary mt-3 text-capitalize"),
									g.Textf("Create %s", levelType.String()),
								),
								Button(
									Type("button"),
									htmx.Get(s.buildRoute("dashboard-timer", "timerID", timerID)),
									Class("btn btn-sm btn-danger ms-2 mt-3 text-capitalize"),
									g.Text("Cancel"),
								),
							),
						),
					),
				),
			),
		),
	)

}

func NewDashboardEditTimerLevelProps(level *poker.TimerLevel, errors []string) *DashboardEditTimerLevelProps {
	return &DashboardEditTimerLevelProps{level, errors}
}

type DashboardEditTimerLevelProps struct {
	level  *poker.TimerLevel
	errors []string
}

func (s *Service) DashboardEditTimerLevelComponent(ctx context.Context, props *DashboardEditTimerLevelProps) g.Node {

	level, errors := props.level, props.errors

	return Div(
		Class("row"),
		Div(
			Class("col"),
			Div(
				Class("card"),
				ID("card-new-level"),
				Div(
					Class("card-header text-center text-capitalize"),
					g.Textf("Edit %s", level.Type.String()),
				),
				Div(
					Class("card-body"),
					s.renderErrorAlert(errors),

					FormEl(
						htmx.Post(s.buildRoute("dashboard-timer-level", "timerID", level.TimerID, "levelID", level.ID)),
						htmx.Target("#card-new-level"),
						Div(
							g.If(
								level.Type == poker.LevelTypeBlind,
								group(
									Class("row row-cols-lg-3 align-items-center"),

									Div(
										Class("col-12"),
										Label(g.Text("Small Blind")),
										Input(
											Class("form-control"), Type("number"), Name("SmallBlind"), Value(format(level.SmallBlind)),
										),
									),
									Div(
										Class("col-12"),
										Label(g.Text("Big Blind")),
										Input(
											Class("form-control"), Type("number"), Name("BigBlind"), Value(format(level.BigBlind)),
										),
									),

									Div(
										Class("col-12"),
										Label(g.Text("Duration (minutes)")),
										Input(
											Class("form-control"), Type("number"), Name("DurationMin"), Value(format(level.DurationMin)),
										),
									),
								),
							),
							g.If(
								level.Type == poker.LevelTypeBreak,
								group(
									Class("row row-cols-lg-3 justify-content-center align-items-center"),

									Div(
										Class("col-12"),
										Label(g.Text("Duration (minutes)")),
										Input(
											Class("form-control"), Type("number"), Name("DurationMin"), Value(format(level.DurationMin)),
										),
									),
								),
							),
						),
						Div(
							Class("row"),
							Div(
								Class("col d-flex justify-content-center"),
								Input(
									Type("hidden"), Name("Type"), Value(level.Type.String()),
								),
								Button(
									Type("submit"),
									Class("btn btn-sm btn-primary mt-3 text-capitalize"),
									g.Textf("Update %s", level.Type.String()),
								),
								Button(
									Type("button"),
									htmx.Get(s.buildRoute("dashboard-timer", "timerID", level.TimerID)),
									Class("btn btn-sm btn-danger ms-2 mt-3 text-capitalize"),
									g.Text("Cancel"),
								),
							),
						),
					),
				),
			),
		),
	)

}

func group(nodes ...g.Node) g.Node {
	return g.Group(nodes)
}
