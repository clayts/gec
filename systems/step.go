package systems

import (
	"time"

	"github.com/clayts/gec/set"
)

type Step struct {
	Components set.Set[func()]
	Time       time.Time
	Duration   time.Duration
	Total      time.Duration
}

func NewStep() *Step {
	u := &Step{}

	u.Time = time.Now()

	return u
}

func (u *Step) Step() {
	now := time.Now()
	u.Duration = now.Sub(u.Time)
	u.Total += u.Duration
	u.Time = now

	u.Components.All(func(e *set.Entity[func()]) bool {
		e.Contents()
		return true
	})
}
