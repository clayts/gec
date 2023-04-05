package random

import (
	"math/rand"

	"github.com/ojrac/opensimplex-go"
)

type Source struct {
	*rand.Rand
	noise [4]opensimplex.Noise
}

func NewSource(seed int64) *Source {
	s := &Source{}
	s.Rand = rand.New(rand.NewSource(seed))
	for i := range s.noise {
		s.noise[i] = opensimplex.NewNormalized(s.Int63())
	}
	return s
}

func (s *Source) Noise1D(x float64) float64 {
	return s.noise[0].Eval4(x, 0, 0, 0)
}

func (s *Source) Noise2D(x, y float64) float64 {
	return s.noise[1].Eval4(x, y, 0, 0)
}

func (s *Source) Noise3D(x, y, z float64) float64 {
	return s.noise[2].Eval4(x, y, z, 0)
}

func (s *Source) Noise4D(x, y, z, w float64) float64 {
	return s.noise[3].Eval4(x, y, z, w)
}
