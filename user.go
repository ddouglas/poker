package poker

import (
	"time"
)

type User struct {
	ID         string
	EmployeeID *uint
	Email      string
	Name       string
	CreatedAt  time.Time
	UpdateAt   time.Time
}
