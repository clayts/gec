package sequence

import (
	"math"
	"sort"
)

type Sequence[T any] struct {
	phases           []T
	cumulativeLength []float64
}

func Make[T any](phases []struct {
	Length   float64
	Contents T
}) Sequence[T] {
	s := Sequence[T]{
		phases:           make([]T, len(phases)),
		cumulativeLength: make([]float64, len(phases)),
	}
	total := 0.0
	for i, phase := range phases {
		total += phase.Length
		s.phases[i] = phase.Contents
		s.cumulativeLength[i] = total
	}
	return s
}

func MakeEqualLength[T any](length float64, contents []T) Sequence[T] {
	s := Sequence[T]{
		phases:           make([]T, len(contents)),
		cumulativeLength: make([]float64, len(contents)),
	}
	total := 0.0
	for i, c := range contents {
		total += length
		s.phases[i] = c
		s.cumulativeLength[i] = total
	}
	return s
}

func (s Sequence[T]) Length() float64 { return s.cumulativeLength[len(s.cumulativeLength)-1] }

func (s Sequence[T]) Get(f float64) T {
	i := sort.SearchFloat64s(s.cumulativeLength, f)
	return s.phases[i]
}

func (s Sequence[T]) GetClamped(f float64) T {
	if f > s.Length() {
		return s.phases[len(s.phases)-1]
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
