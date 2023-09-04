package poker

import (
	"time"
)

type Timer struct {
	ID           string `schema:"-"`
	UserID       string `schema:"-"`
	Name         string
	Levels       []*TimerLevel
	CurrentLevel uint      `schema:"-"`
	CreatedAt    time.Time `schema:"-"`
	UpdatedAt    time.Time `schema:"-"`
}

type TimerType string

const (
	TimerTypeBlind TimerType = "blind"
	TimerTypeBreak TimerType = "break"
)

func (tt TimerType) String() string {
	return string(tt)
}

var allTimerTypes = []TimerType{TimerTypeBlind, TimerTypeBreak}

func (tt TimerType) Valid() bool {
	for _, t := range allTimerTypes {
		if t == tt {
			return true
		}
	}
	return false
}

type TimerLevel struct {
	ID          string
	Type        TimerType
	TimerID     string
	Level       float64
	SmallBlind  float64
	BigBlind    float64
	Ante        float64
	DurationMin float64
	DurationSec float64
	// DurationStr is a the string representation of the remaining time in the MM:SS format
	DurationStr string
	ShowHours   bool
}
