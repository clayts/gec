package procgen

import "github.com/ojrac/opensimplex-go"

type Source4D func(x, y, z, w float64) float64
type Source3D func(x, y, z float64) float64
type Source2D func(x, y float64) float64

func (s4 Source4D) Slice(w float64) Source3D {
	return func(x, y, z float64) float64 {
		return s4(x, y, z, w)
	}
}

func (s3 Source3D) Slice(z float64) Source2D {
	return func(x, y float64) float64 {
		return s3(x, y, z)
	}
}

func NewNoise(seed int64) Source4D {
	return opensimplex.NewNormalized(seed).Eval4
}
