package systems

import (
	"time"

	"github.com/clayts/gec/set"
)

type Update struct {
	Components   set.Set[func()]
	StepTime     time.Time
	StepDuration time.Duration
}

func NewUpdate() *Update {
	u := &Update{}

	u.StepTime = time.Now()

	return u
}

func (u *Update) Step() {
	now := time.Now()
	u.StepDuration = now.Sub(u.StepTime)
	u.StepTime = now

	u.Components.All(func(e *set.Entity[func()]) bool {
		e.Contents()
		return true
	})
}
