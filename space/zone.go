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

func (z *Zone[A]) SetState(state bool) *Zone[A] {
	if state {
		if z.subSpace == nil {
			z.space.add(z)
		}
	} else {
		if z.subSpace != nil {
			z.space.remove(z)
		}
	}
	return z
}

func (z *Zone[A]) State() bool { return z.subSpace != nil }

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
