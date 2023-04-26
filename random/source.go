package random

import (
	"math/rand"

	"github.com/ojrac/opensimplex-go"
)

type Source struct {
	*rand.Rand
	Noise func(coordinates ...float64) float64
}

func NewSource(seed int64) *Source {
	s := &Source{}
	s.Rand = rand.New(rand.NewSource(seed))

	// This hack is required because implementations of openSimplex noise in different dimensions produce a different distribution of outputs.
	noiseA := opensimplex.NewNormalized(s.Rand.Int63()).Eval4
	noiseB := opensimplex.NewNormalized(s.Rand.Int63()).Eval4
	noiseC := opensimplex.NewNormalized(s.Rand.Int63()).Eval4
	noiseD := opensimplex.NewNormalized(s.Rand.Int63()).Eval4
	noiseE := opensimplex.NewNormalized(s.Rand.Int63()).Eval4

	s.Noise = func(coordinates ...float64) float64 {
		switch len(coordinates) {
		case 0:
			return noiseA(s.Float64()*1024, 0, 0, 0)
		case 1:
			return noiseB(coordinates[0], 0, 0, 0)
		case 2:
			return noiseC(coordinates[0], coordinates[1], 0, 0)
		case 3:
			return noiseD(coordinates[0], coordinates[1], coordinates[2], 0)
		case 4:
			return noiseE(coordinates[0], coordinates[1], coordinates[2], coordinates[3])
		default:
			panic("noise coordinates must be 1, 2, 3 or 4 dimensional")
		}
	}

	return s
}
