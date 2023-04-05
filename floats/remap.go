package floats

import "math"

func Remap(value, fromMin, fromMax, toMin, toMax float64) float64 {
	return (value * ((toMax - toMin) / (fromMax - fromMin))) + (toMin - fromMin)
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

func Mirror(value, min, max float64) float64 {
	max2 := max + (max - min)
	w := Wrap(value, min, max2)
	if w > max {
		return Remap(w, max, max2, max, min)
	}
	return w
}
