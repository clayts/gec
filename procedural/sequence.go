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

func Make[T any](steps ...struct {
	Length float64
	Step   T
}) Sequence[T] {
	s := Sequence[T]{
		steps:            make([]T, len(steps)),
		cumulativeLength: make([]float64, len(steps)),
	}
	total := 0.0
	for i, step := range steps {
		total += step.Length
		s.steps[i] = step.Step
		s.cumulativeLength[i] = total
	}
	return s
}

func MakeIndexed[T any](steps ...T) Sequence[T] {
	return Sequence[T]{steps: steps}
}

func (s Sequence[T]) Length() float64 {
	if s.cumulativeLength == nil {
		return float64(len(s.cumulativeLength))
	}
	return s.cumulativeLength[len(s.cumulativeLength)-1]
}

func (s Sequence[T]) Get(f float64) T {
	if s.cumulativeLength == nil {
		return s.steps[int(f)]
	}
	i := sort.SearchFloat64s(s.cumulativeLength, f)
	return s.steps[i]
}

func (s Sequence[T]) GetClamped(f float64) T {
	if f > s.Length() {
		return s.steps[len(s.steps)-1]
	}
	return s.Get(f)
}

func (s Sequence[T]) GetWrapped(f float64) T {
	return s.Get(math.Remainder(f, s.Length()))
}

func (s Sequence[T]) GetMirrored(f float64) T {
	w := s.Length() * 2
	r := math.Remainder(f, w)
	if r < s.Length() {
		return s.Get(f)
	}
	return s.Get(w - f)
}

func (s Sequence[T]) GetFraction(f float64) T {
	return s.Get(f * s.Length())
}

func (s Sequence[T]) GetRandom() T {
	return s.GetFraction(rand.Float64())
}
