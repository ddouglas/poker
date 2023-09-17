package poker

import (
	"fmt"
	"strings"
	"time"
)

type Timer struct {
	ID           string `schema:"-"`
	UserID       string `schema:"-"`
	Name         string
	Levels       []*TimerLevel
	CurrentLevel uint      `schema:"-"`
	IsComplete   bool      `schema:"-"`
	CreatedAt    time.Time `schema:"-"`
	UpdatedAt    time.Time `schema:"-"`
}

func (t Timer) Validate() error {

	if t.ID == "" {
		return fmt.Errorf("id cannot be empty")
	}

	if t.UserID == "" {
		return fmt.Errorf("user id cannot be empty")
	}

	if len(t.Name) < 3 {
		return fmt.Errorf("name must be 3 or more characters in length")
	}

	return nil

}

type LevelType string

const (
	LevelTypeBlind LevelType = "blind"
	LevelTypeBreak LevelType = "break"
)

func (tt LevelType) String() string {
	return string(tt)
}

var allLevelTypes = []LevelType{LevelTypeBlind, LevelTypeBreak}

func (tt LevelType) Valid() bool {
	for _, t := range allLevelTypes {
		if t == tt {
			return true
		}
	}
	return false
}

var strAllLevelTypes = []string{LevelTypeBlind.String(), LevelTypeBreak.String()}

type TimerLevel struct {
	ID          string
	Type        LevelType
	TimerID     string
	Level       float64
	SmallBlind  float64
	BigBlind    float64
	Ante        float64
	DurationMin float64
	DurationSec float64

	// DurationStr is a the string representation of the remaining time in the MM:SS format
	DurationStr string
}

func (t TimerLevel) Validate() error {

	if t.ID == "" {
		return fmt.Errorf("id cannot be empty")
	}

	if !t.Type.Valid() {
		return fmt.Errorf("type is not a valid type, expected one of: %s", strings.Join(strAllLevelTypes, ","))
	}

	if t.TimerID == "" {
		return fmt.Errorf("timer id cannot be empty")
	}

	if t.Level < 0 {
		return fmt.Errorf("level must be greater than or equal to 0")
	}

	if t.Type == LevelTypeBlind {
		if t.SmallBlind < 0 {
			return fmt.Errorf("small blind must be greater than or equal to 0")
		}

		if t.BigBlind < 0 {
			return fmt.Errorf("large blind must be greater than or equal to 0")
		}

		if t.SmallBlind > t.BigBlind {
			return fmt.Errorf("small blind cannot be greater than big blind")
		}
	}

	if t.DurationMin <= 0 {
		return fmt.Errorf("duration must be greater than 0")
	}

	return nil

}

func (t TimerLevel) AudioS3Key() string {
	switch t.Type {
	case LevelTypeBlind:
		return fmt.Sprintf("%.0f-%.0f-%.0f", t.SmallBlind, t.BigBlind, t.DurationMin)
	case LevelTypeBreak:
		return fmt.Sprintf("%.0f", t.DurationMin)
	}

	return "unrecognized-level-type"
}
