package geometry

import "math"

type Angle float64

func (a Angle) Radians() float64 {
	return float64(a) * math.Pi / 180
}

func (a Angle) SinCos() (sin, cos float64) {
	return math.Sincos(a.Radians())
}
