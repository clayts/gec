package space

import geo "github.com/clayts/gec/geometry"

type Zone[A any] struct {
	shape    geo.Shape
	space    *Space[A]
	subSpace *subSpace[A]
	index    int
	Contents A
}

func (spc *Space[A]) New(s geo.Shape) *Zone[A] {
	return &Zone[A]{space: spc, shape: s}
}

func (z *Zone[A]) Enable() *Zone[A] {
	if z.subSpace == nil {
		z.space.add(z)
	}
	return z
}

func (z *Zone[A]) Disable() *Zone[A] {
	if z.subSpace != nil {
		z.space.remove(z)
	}
	return z
}

func (z *Zone[A]) SetShape(s geo.Shape) *Zone[A] {
	z.shape = s
	if z.subSpace != nil {
		z.space.remove(z)
		z.space.add(z)
	}
	return z
}

func (z *Zone[A]) SetContents(c A) *Zone[A] {
	z.Contents = c
	return z
}

func (z *Zone[A]) Shape() geo.Shape { return z.shape }
