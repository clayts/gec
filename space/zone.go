package space

import geo "github.com/clayts/gec/geometry"

type Zone[T any] struct {
	shape    geo.Shape
	space    *Space[T]
	subSpace *subSpace[T]
	index    int
	Contents T
}

func (spc *Space[T]) NewZone() *Zone[T] {
	return &Zone[T]{space: spc}
}

func (z *Zone[T]) Enable() *Zone[T] {
	if z.subSpace == nil {
		z.space.add(z)
	}
	return z
}

func (z *Zone[T]) Disable() *Zone[T] {
	if z.subSpace != nil {
		z.space.remove(z)
	}
	return z
}

func (z *Zone[T]) SetShape(s geo.Shape) *Zone[T] {
	z.shape = s
	if z.subSpace != nil {
		z.space.remove(z)
		z.space.add(z)
	}
	return z
}

func (z *Zone[T]) SetContents(c T) *Zone[T] {
	z.Contents = c
	return z
}

func (z *Zone[T]) Shape() geo.Shape { return z.shape }
