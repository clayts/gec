package random

type NoiseSelector[A any] func(coordinates ...float64) A

func MakeNoiseSelector[A any](n Noise, s Selector[A]) NoiseSelector[A] {
	return func(coordinates ...float64) A {
		return s(n(coordinates...))
	}
}
