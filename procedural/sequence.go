package procedural

import (
	"math"
	"math/rand"
	"sort"
)

type Sequence[T any] struct {
	steps            []T
	cumulativeLength []float64
}

func NewSequence[T any]() *Sequence[T] {
	s := &Sequence[T]{}

	return s
}

func (s *Sequence[T]) Add(length float64, step T) *Sequence[T] {
	s.steps = append(s.steps, step)
	if len(s.cumulativeLength) == 0 {
		s.cumulativeLength = []float64{length}
	} else {
		prev := s.cumulativeLength[len(s.cumulativeLength)-1]
		s.cumulativeLength = append(s.cumulativeLength, length+prev)
	}
	return s
}

func (s Sequence[T]) Length() float64 {
	return s.cumulativeLength[len(s.cumulativeLength)-1]
}

func (s Sequence[T]) Get(f float64) T {
	i := sort.SearchFloat64s(s.cumulativeLength, f)
	return s.steps[i]
}

func (s Sequence[T]) GetClamped(f float64) T {
	if f > s.Length() {
		return s.steps[len(s.steps)-1]
	}
	if f < 0 {
		return s.steps[0]
	}
	return s.Get(f)
}

func (s Sequence[T]) GetWrapped(f float64) T {
	return s.Get(math.Remainder(f, s.Length()))
}

func (s Sequence[T]) GetMirrored(f float64) T {
	return s.Get(math.Abs(math.Remainder(f, s.Length()*2)))
}

func (s Sequence[T]) GetFraction(f float64) T {
	return s.Get(f * s.Length())
}

func (s Sequence[T]) GetRandom() T {
	return s.GetFraction(rand.Float64())
}
