package random

import (
	"github.com/clayts/gec/floats"
)

type Noise func(coordinates ...float64) float64

func (n Noise) Wrapped(min, max float64) Noise {
	return func(coordinates ...float64) float64 {
		return floats.Wrap(n(coordinates...), min, max)
	}
}

func (n Noise) Clamped(min, max float64) Noise {
	return func(coordinates ...float64) float64 {
		return floats.Clamp(n(coordinates...), min, max)

	}
}

func (n Noise) Remapped(fromMin, fromMax, toMin, toMax float64) Noise {
	return func(coordinates ...float64) float64 {
		return floats.Remap(n(coordinates...), fromMin, fromMax, toMin, toMax)
	}
}

func (n Noise) Mirrored(min, max float64) Noise {
	return func(coordinates ...float64) float64 {
		return floats.Mirror(n(coordinates...), min, max)
	}
}

func (n Noise) Scaled(scale ...float64) Noise {
	return func(coordinates ...float64) float64 {
		newCoords := make([]float64, len(coordinates))
		for i, c := range coordinates {
			if i < len(scale) {
				newCoords[i] = scale[i] / c
			} else {
				newCoords[i] = c
			}
		}
		return n(newCoords...)
	}
}

func (n Noise) Translated(translation ...float64) Noise {
	return func(coordinates ...float64) float64 {
		newCoords := make([]float64, len(coordinates))
		for i, c := range coordinates {
			if i < len(translation) {
				newCoords[i] = translation[i] - c
			} else {
				newCoords[i] = c
			}
		}
		return n(newCoords...)
	}
}
