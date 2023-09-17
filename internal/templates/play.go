package templates

import (
	"context"
	"fmt"
	"poker"
	"time"

	g "github.com/maragudk/gomponents"
	htmx "github.com/maragudk/gomponents-htmx"
	. "github.com/maragudk/gomponents/html"
)

type PlayProps struct {
	User         *poker.User
	Timer        *poker.Timer
	Level        *poker.TimerLevel
	CurrentLevel uint
}

func (s *Service) Play(ctx context.Context, props *PlayProps) g.Node {

	return Doctype(
		HTML(
			Lang("en"),
			s.gtop(ctx),
			Body(
				s.gnavbar(ctx),
				s.TimerMasthead(ctx, props.Timer, props.Level),
				s.gbottom(),
				Script(
					Src(fmt.Sprintf("%s/js/countdown.js?v=%d", s.buildRoute("static"), time.Now().Unix())),
				),
			),
		),
	)

}

func (s *Service) TimerMasthead(ctx context.Context, timer *poker.Timer, level *poker.TimerLevel) g.Node {

	var nextLevel *poker.TimerLevel = nil
	if int(timer.CurrentLevel+1) <= len(timer.Levels)-1 {
		nextLevel = timer.Levels[timer.CurrentLevel+1]
	}

	return Div(
		ID("timer-container"), Class("container"), htmx.SwapOOB("true"),
		Div(
			Audio(
				ID("audio-play"),
				Source(
					Src(s.buildRoute("dashboard-timer-level-audio", "timerID", level.TimerID, "levelID", level.ID, "action", "play")),
					Type("audio/mpeg"),
				),
				// DataAttr("continue-audio", s.buildRoute("dashboard-timer-level-audio", "timerID", level.TimerID, "levelID", level.ID, "action", "contiue")),
			),
			Audio(
				ID("audio-continue"),
				Source(
					Src(s.buildRoute("dashboard-timer-level-audio", "timerID", level.TimerID, "levelID", level.ID, "action", "continue")),
					Type("audio/mpeg"),
				),
			),
			Audio(
				ID("audio-beep"),
				Source(
					Src("/static/audio/10_sec_beep_countdown.mp3"),
					Type("audio/mpeg"),
				),
			),
		),
		Div(
			Class("row"),
			Div(
				Class("col-10 offset-1"),
				H1(
					Class("text-center"),
					g.Textf("Timer %s", timer.Name),
				),
				Hr(),
				Div(
					Class("row"),
					Div(
						Class("col"),
						Div(
							Class("timer-container d-flex justify-content-center align-items-center"),
							g.If(
								timer.IsComplete,
								Div(
									ID("timer"), Class("timer-complete-font"),
									g.Text("Timer Complete"),
								),
							),
							g.If(
								!timer.IsComplete,
								Div(
									ID("timer"), Class("timer-large-font"), DataAttr("level-duration-sec", fmt.Sprintf("%v", level.DurationSec)),
									g.Text(level.DurationStr),
								),
							),
						),
					),
				),
				Div(
					Class("row mt-2"),
					Div(
						Class("col-4"),
						Div(
							Class("d-flex justify-content-center"),
							s.formatPlayLevelDisplay(ctx, "Current Blind", level),
						),
					),
					Div(
						Class("col-4"),
						Div(
							Class("row"),

							Div(
								Class("row mt-2"),
								s.formatTimerButtons(ctx, timer, level),
							),
						),
					),
					Div(
						Class("col-4"),
						Div(
							Class("d-flex justify-content-center"),
							s.formatPlayLevelDisplay(ctx, "Next Blind", nextLevel),
						),
					),
				),
			),
		),
	)
}

func (s *Service) formatTimerButtons(ctx context.Context, timer *poker.Timer, level *poker.TimerLevel) g.Node {

	nodes := make([]g.Node, 0)
	if timer.CurrentLevel > 0 {
		nodes = append(
			nodes,
			Div(
				Class("col text-center"),
				I(
					ID("trigger-previous-timer-level"),
					Class("fa-solid fa-angles-left fa-3x"),
					htmx.Get(s.buildRoute("play-timer-previous-level", "timerID", level.TimerID)),
				),
			),
		)
	}

	if timer.IsComplete {
		nodes = append(
			nodes,
			Div(
				Class("col text-center"),
				I(
					ID("trigger-reset-timer-level"),
					Class("fa-solid fa-arrow-rotate-left fa-3x"),
					htmx.Get(s.buildRoute("play-timer-reset-level", "timerID", level.TimerID)),
				),
			),
		)
	}

	nodes = append(
		nodes,
		Div(
			Class("col text-center"),
			I(
				ID("toggle-timer-button"),
				Class("fa-solid fa-circle-play fa-3x"),
			),
		),
	)

	if int(timer.CurrentLevel) < len(timer.Levels)-1 {
		nodes = append(
			nodes,
			Div(
				Class("col text-center"),
				I(
					ID("trigger-next-timer-level"),
					Class("fa-solid fa-angles-right fa-3x"),
					htmx.Get(s.buildRoute("play-timer-next-level", "timerID", level.TimerID)),
				),
			),
		)
	}

	return g.Group(nodes)

}

func (s *Service) formatPlayLevelDisplay(ctx context.Context, header string, level *poker.TimerLevel) g.Node {

	nodes := make([]g.Node, 0)
	nodes = append(
		nodes,
		Class("text-center"),
		g.Text(header),
		Hr(),
	)

	if level != nil {
		nodes = append(
			nodes,
			group(
				g.If(
					level.Type == poker.LevelTypeBlind,
					g.Textf("%.0f / %.0f", level.SmallBlind, level.BigBlind),
				),
				g.If(
					level.Type == poker.LevelTypeBreak,
					g.Text("Break"),
				),
			),
		)
	} else if level == nil {
		nodes = append(nodes, g.Text("No More Blinds"))
	}

	return H1(nodes...)

}
