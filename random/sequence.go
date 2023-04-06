package random

import (
	"fmt"
	"sort"

	"github.com/clayts/gec/floats"
)

type Sequence[A any] []struct {
	Item   A
	Length float64
}

func MakeSequence[A any](slice ...A) Sequence[A] {
	s := make(Sequence[A], len(slice))
	for i, item := range slice {
		s[i].Item = item
		s[i].Length = 1
	}
	return s
}

func (s Sequence[A]) analyze() ([]A, []float64, float64, bool) {

	homogenous := true
	firstLength := s[0].Length

	items := make([]A, len(s))
	items[0] = s[0].Item

	cumulative := make([]float64, len(s))
	cumulative[0] = firstLength

	total := firstLength
	for i := 1; i < len(s); i++ {
		phase := s[i]
		if phase.Length != firstLength {
			homogenous = false
		}
		total += phase.Length
		items[i] = phase.Item
		cumulative[i] = total
	}
	return items, cumulative, total, homogenous
}

func (s Sequence[A]) Selector() func(index float64) A {
	items, cumulative, total, homogenous := s.analyze()
	if homogenous {
		return func(index float64) A {
			if index < 0 || index > 1 {
				panic(fmt.Sprint("index out of bounds", index))
			}
			if index == 1 {
				return items[len(items)-1]
			}
			return items[int(floats.Remap(index, 0, 1, 0, float64(len(items))))]
		}
	}
	return func(index float64) A {
		if index < 0 || index > 1 {
			panic(fmt.Sprint("index out of bounds", index))
		}
		return items[sort.SearchFloat64s(cumulative, floats.Remap(index, 0, 1, 0, total))]
	}
}
