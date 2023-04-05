package random

import (
	"fmt"
	"sort"

	"github.com/clayts/gec/floats"
)

type Selector[A any] func(index float64) A

type Phase[A any] struct {
	Item   A
	Length float64
}

type Sequence[A any] []Phase[A]

func MakeSequence[A any](slice ...A) Sequence[A] {
	s := make(Sequence[A], len(slice))
	for i, item := range slice {
		s[i].Item = item
		s[i].Length = 1
	}
	return s
}

func (s Sequence[A]) compile() ([]A, []float64, float64) {
	items := make([]A, len(s))
	cumulative := make([]float64, len(s))
	total := 0.0
	for i, phase := range s {
		total += phase.Length
		items[i] = phase.Item
		cumulative[i] = total
	}
	return items, cumulative, total
}

func (s Sequence[A]) Length() float64 {
	total := 0.0
	for _, phase := range s {
		total += phase.Length
	}
	return total
}

func (s Sequence[A]) Selector(min, max float64) Selector[A] {
	items, cumulative, total := s.compile()
	return func(index float64) A {
		if index < min || index > max {
			panic(fmt.Sprint("index out of bounds", index, "(min:", min, "max:", max, ")"))
		}
		return items[sort.SearchFloat64s(cumulative, floats.Remap(index, min, max, 0, total))]
	}
}

func MakeSelector[A any](min, max float64, slice ...A) Selector[A] {
	return func(index float64) A {
		if index < min || index > max {
			panic(fmt.Sprint("index out of bounds", index, "(min:", min, "max:", max, ")"))
		}
		if index == max {
			return slice[len(slice)-1]
		}
		return slice[int(floats.Remap(index, min, max, 0, float64(len(slice))))]
	}
}
