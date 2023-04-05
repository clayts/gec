package random

import (
	"math"
	"sort"
)

type Sequence[T any] []struct {
	Item   T
	Length float64
}

func MakeSequence[T any](slice ...T) Sequence[T] {
	s := make(Sequence[T], len(slice))
	for i, item := range slice {
		s[i].Item = item
		s[i].Length = 1
	}
	return s
}

func (s Sequence[T]) compile() ([]T, []float64, float64) {
	items := make([]T, len(s))
	cumulative := make([]float64, len(s))
	total := 0.0
	for i, phase := range s {
		total += phase.Length
		items[i] = phase.Item
		cumulative[i] = total
	}
	return items, cumulative, total
}

func (s Sequence[T]) Length() float64 {
	total := 0.0
	for _, phase := range s {
		total += phase.Length
	}
	return total
}

func (s Sequence[T]) Selector(min, max float64) func(time float64) T {
	items, cumulative, total := s.compile()
	return func(time float64) T {
		return items[sort.SearchFloat64s(cumulative, Remap(time, min, max, 0, total))]
	}
}

func MakeSelector[T any](min, max float64, slice ...T) func(time float64) T {
	return func(time float64) T {
		if time == max {
			return slice[len(slice)-1]
		}
		return slice[int(Remap(time, min, max, 0, float64(len(slice))))]
	}
}

func Wrap(value, min, max float64) float64 {
	m := math.Mod(value, max-min)
	if m == 0 {
		return min
	}
	shift := min
	if value < 0 {
		shift = max
	}
	return m + shift
}

func Clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func Remap(value, fromMin, fromMax, toMin, toMax float64) float64 {
	return (value * ((toMax - toMin) / (fromMax - fromMin))) + (toMin - fromMin)
}

// func Continuum[T any](items []T) func(f float64) T {
// 	l := float64(len(items))
// 	return func(f float64) T {
// 		i := int(f * l)
// 		return items[i]
// 	}
// }

// type Bias[T any] struct {
// 	Bias float64
// 	Item T
// }

// func Biases[T any](items []T) []Bias[T] {
// 	bs := make([]Bias[T], len(items))
// 	for i, item := range items {
// 		bs[i] = Bias[T]{Bias: 1, Item: item}
// 	}
// 	return bs
// }

// func BiasedContinuum[T any](items []Bias[T]) func(f float64) T {
// 	cumulative := make([]float64, len(items))
// 	for i, b := range items {
// 		prev := 0.0
// 		if i > 0 {
// 			prev = cumulative[i-1]
// 		}
// 		cumulative[i] = prev + b.Bias
// 	}
// 	return func(f float64) T {
// 		i := sort.SearchFloat64s(cumulative, f)
// 		return items[i].Item
// 	}
// }

// type Sequence[T any] struct {
// 	steps            []T
// 	cumulativeLength []float64
// }

// func NewSequence[T any]() *Sequence[T] {
// 	s := &Sequence[T]{}

// 	return s
// }

// func (s *Sequence[T]) Add(length float64, step T) *Sequence[T] {
// 	s.steps = append(s.steps, step)
// 	if len(s.cumulativeLength) == 0 {
// 		s.cumulativeLength = []float64{length}
// 	} else {
// 		prev := s.cumulativeLength[len(s.cumulativeLength)-1]
// 		s.cumulativeLength = append(s.cumulativeLength, length+prev)
// 	}
// 	return s
// }

// func (s *Sequence[T]) Length() float64 {
// 	return s.cumulativeLength[len(s.cumulativeLength)-1]
// }

// func (s *Sequence[T]) Get(f float64) T {
// 	i := sort.SearchFloat64s(s.cumulativeLength, f)
// 	return s.steps[i]
// }

// func (s *Sequence[T]) GetClamped(f float64) T {
// 	if f > s.Length() {
// 		return s.steps[len(s.steps)-1]
// 	}
// 	if f < 0 {
// 		return s.steps[0]
// 	}
// 	return s.Get(f)
// }

// func (s *Sequence[T]) GetWrapped(f float64) T {
// 	return s.Get(math.Remainder(f, s.Length()))
// }

// func (s *Sequence[T]) GetMirrored(f float64) T {
// 	return s.Get(math.Abs(math.Remainder(f, s.Length()*2)))
// }

// func (s *Sequence[T]) GetFraction(f float64) T {
// 	return s.Get(f * s.Length())
// }

// func (s *Sequence[T]) GetRandom() T {
// 	return s.GetFraction(rand.Float64())
// }
