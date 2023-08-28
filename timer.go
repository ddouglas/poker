package poker

import (
	"time"
)

type Timer struct {
	ID        string `schema:"-"`
	UserID    string `schema:"-"`
	Name      string
	Levels    []*TimerLevel
	CreatedAt time.Time `schema:"-"`
	UpdatedAt time.Time `schema:"-"`
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
	Level       uint
	SmallBlind  uint
	BigBlind    uint
	Ante        uint
	DurationSec uint
	TSCreated   time.Time
	TSUpdated   time.Time
}
