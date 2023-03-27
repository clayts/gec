package procgen

import (
	"math/rand"
)

func GetRandom[T any](s []T) T {
	return s[rand.Intn(len(s))]
}

func GetFraction[T any](f float64, s []T) T {
	v := f * float64(len(s))
	i := int(v)
	return s[i]
}
