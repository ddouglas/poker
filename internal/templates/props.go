package templates

import (
	"context"
	"io"
	"poker"
	"poker/internal"
)

type UserProps interface {
	User() *poker.User
}

type NavbarProps interface {
	ContextProps
	UserProps
}

type ContextProps interface {
	Context() context.Context
}

type TimersProps interface {
	Timers() []*poker.Timer
}

type TimerProps interface {
	Timer() *poker.Timer
}

type HomepageProps struct {
	ctx context.Context
}

func (hp HomepageProps) Context() context.Context {
	return hp.ctx
}

func (hp HomepageProps) User() *poker.User {
	return internal.UserFromContext(hp.ctx)
}

type DashboardProps interface {
	ContextProps
	UserProps
	io.Writer
}

func NewDashboardProps(ctx context.Context, w io.Writer) *dashboardProps {
	return &dashboardProps{
		ctx:    ctx,
		Writer: w,
	}
}

type dashboardProps struct {
	ctx context.Context
	io.Writer
}

func (hp dashboardProps) Context() context.Context {
	return hp.ctx
}

func (hp dashboardProps) User() *poker.User {
	return internal.UserFromContext(hp.ctx)
}

type DashboardTimersProps interface {
	ContextProps
	TimersProps
	UserProps
	io.Writer
}

func NewDashboardTimersProps(ctx context.Context, timers []*poker.Timer, w io.Writer) DashboardTimersProps {
	return &dashboardTimersProps{
		ctx:    ctx,
		timers: timers,
		Writer: w,
	}
}

type dashboardTimersProps struct {
	ctx    context.Context
	timers []*poker.Timer
	io.Writer
}

func (dtp dashboardTimersProps) Context() context.Context {
	return dtp.ctx
}

func (dtp dashboardTimersProps) Timers() []*poker.Timer {
	return dtp.timers
}

func (hp dashboardTimersProps) User() *poker.User {
	return internal.UserFromContext(hp.ctx)
}

type DashboardTimerProps interface {
	ContextProps
	UserProps
	TimerProps
	io.Writer
}

func NewDashboardTimerProps(ctx context.Context, timer *poker.Timer, w io.Writer) DashboardTimerProps {
	return &dashboardTimerProps{
		ctx:    ctx,
		timer:  timer,
		Writer: w,
	}
}

type dashboardTimerProps struct {
	ctx   context.Context
	timer *poker.Timer
	io.Writer
}

func (dtp dashboardTimerProps) Context() context.Context {
	return dtp.ctx
}
func (dtp dashboardTimerProps) User() *poker.User {
	return internal.UserFromContext(dtp.ctx)
}

func (dtp dashboardTimerProps) Timer() *poker.Timer {
	println("(dtp dashboardTimerProps) Timer() *poker.Timer")
	return dtp.timer
}

type TimerLevelProps interface {
	ContextProps
	TimerProps
}

func NewTimerLevelProps(ctx context.Context, timer *poker.Timer) TimerLevelProps {
	return &timerLevelProps{ctx, timer}
}

type timerLevelProps struct {
	ctx   context.Context
	timer *poker.Timer
}

func (dtp timerLevelProps) Context() context.Context {
	return dtp.ctx
}

func (dtp timerLevelProps) Timer() *poker.Timer {
	return dtp.timer
}

type TimerLevelEditProps interface {
	ContextProps
	Level() *poker.TimerLevel
	LevelIndex() int
}

func NewTimerLevelEditProps(ctx context.Context, timer *poker.Timer, idx int) TimerLevelEditProps {
	return &timerLevelEditProps{ctx, timer, idx}
}

type timerLevelEditProps struct {
	ctx   context.Context
	timer *poker.Timer
	idx   int
}

func (tle timerLevelEditProps) Context() context.Context {
	return tle.ctx
}

func (tle timerLevelEditProps) Timer() *poker.Timer {
	return tle.timer
}

func (tle timerLevelEditProps) Level() *poker.TimerLevel {
	return tle.timer.Levels[tle.idx]
}

func (tle timerLevelEditProps) LevelIndex() int {
	return tle.idx
}
