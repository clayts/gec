package space

import (
	geo "github.com/clayts/gec/geometry"
)

type subSpace[A any] struct {
	parent         *subSpace[A]
	center, radius geo.Vector
	children       [2]*subSpace[A]
	vertical       bool
	upper          bool
	zones          []*Zone[A]
	size           int
}

func (b *subSpace[A]) all(f func(z *Zone[A]) bool) bool {
	if b == nil {
		return true
	}

	for i := len(b.zones) - 1; i >= 0; i-- {
		z := b.zones[i]
		if !f(z) {
			return false
		}
	}
	return b.children[0].all(f) && b.children[1].all(f)
}

func (b *subSpace[A]) allIntersectingSkipBounds(s geo.Shape, f func(z *Zone[A]) bool) bool {
	if b == nil {
		return true
	}
	for i := len(b.zones) - 1; i >= 0; i-- {
		z := b.zones[i]
		if (s.ShapeType() == geo.RECTANGLE || geo.AllEdges(s, func(g geo.Segment) bool { return g.AxisAligned() || !geo.AllAnticlockwise(g, z.shape) })) &&
			(z.shape.ShapeType() == geo.RECTANGLE || geo.AllEdges(z.shape, func(g geo.Segment) bool { return g.AxisAligned() || !geo.AllAnticlockwise(g, s) })) {
			if !f(z) {
				return false
			}
		}
	}
	return b.children[0].allIntersectingSkipBounds(s, f) && b.children[1].allIntersectingSkipBounds(s, f)
}

func (b *subSpace[A]) allIntersecting(s geo.Shape, f func(z *Zone[A]) bool) bool {
	if b == nil {
		return true
	}
	if !s.Bounds().Intersects(b.bounds()) {
		return true
	}
	if s.Bounds().Contains(b.bounds()) {
		return b.allIntersectingSkipBounds(s, f)
	}
	for i := len(b.zones) - 1; i >= 0; i-- {
		z := b.zones[i]
		if geo.Intersects(s, z.shape) {
			if !f(z) {
				return false
			}
		}
	}
	return b.children[0].allIntersecting(s, f) && b.children[1].allIntersecting(s, f)
}

func (b *subSpace[A]) remove(z *Zone[A]) {

	if z.index != -1 {
		finalIndex := len(b.zones) - 1
		if finalIndex != 0 {
			finalEntity := b.zones[finalIndex]
			b.zones[z.index] = finalEntity
			finalEntity.index = z.index
			b.zones[finalIndex] = nil
		}
		b.zones = b.zones[:finalIndex]
		z.index = -1
		z.subSpace = nil
	}

	b.decrease()
}

func (b *subSpace[A]) clean() {
	if b.size <= maxCached {
		if b.children[0] != nil {
			for _, z := range b.children[0].zones {
				b.store(z)
			}
			b.children[0] = nil
		}
		if b.children[1] != nil {
			for _, z := range b.children[1].zones {
				b.store(z)
			}
			b.children[1] = nil
		}
	}
}

func (b *subSpace[A]) decrease() {
	b.size--
	b.clean()
	if b.parent != nil {
		b.parent.decrease()
	}
}

func (b *subSpace[A]) add(z *Zone[A]) {
	// if size < max, cache this
	// if size == max, distribute all
	// if size > max,  distribute this
	if b.size < maxCached {
		b.store(z)
	} else {
		if b.size == maxCached {
			leaves := b.zones
			b.zones = b.zones[:0]
			for _, l2 := range leaves {
				b.distribute(l2)
			}
		}
		b.distribute(z)
	}
	b.size++
}

func (b *subSpace[A]) distribute(z *Zone[A]) {
	if !b.vertical {
		if z.shape.Bounds().Max.X < b.center.X {
			b.demandChild(false).add(z)
			return
		} else if z.shape.Bounds().Min.X > b.center.X {
			b.demandChild(true).add(z)
			return
		}
	} else {
		if z.shape.Bounds().Max.Y < b.center.Y {
			b.demandChild(false).add(z)
			return
		} else if z.shape.Bounds().Min.Y > b.center.Y {
			b.demandChild(true).add(z)
			return
		}
	}
	b.store(z)
}

func (b *subSpace[A]) store(z *Zone[A]) {
	z.index = len(b.zones)
	z.subSpace = b
	b.zones = append(b.zones, z)
}

func (b *subSpace[A]) bounds() geo.Rectangle {
	return geo.R(b.center.Minus(b.radius), b.center.Plus(b.radius))
}

func (b *subSpace[A]) demandParent() *subSpace[A] {
	if b.parent == nil {
		if !b.vertical {
			if !b.upper {
				// horizontal, lower (is the lower half of a vertical) - make vertical, upper (will be upper half of a horizontal)
				b.parent = &subSpace[A]{
					center:   b.center.Plus(geo.V(0, b.radius.Y)),
					radius:   b.radius.Times(geo.V(1, 2)),
					vertical: true,
					upper:    true,
					children: [2]*subSpace[A]{b, nil},
					size:     b.size,
				}
			} else {
				// horizontal, upper (is the upper half of a vertical) - make vertical, lower (will be lower half of a horizontal)
				b.parent = &subSpace[A]{
					center:   b.center.Minus(geo.V(0, b.radius.Y)),
					radius:   b.radius.Times(geo.V(1, 2)),
					vertical: true,
					upper:    false,
					children: [2]*subSpace[A]{nil, b},
					size:     b.size,
				}
			}
		} else {
			if !b.upper {
				// vertical, lower  (is the lower half of a horizontal) - make horizontal, lower (will be lower half of a vertical)
				b.parent = &subSpace[A]{
					center:   b.center.Plus(geo.V(b.radius.X, 0)),
					radius:   b.radius.Times(geo.V(2, 1)),
					vertical: false,
					upper:    false,
					children: [2]*subSpace[A]{b, nil},
					size:     b.size,
				}
			} else {
				// vertical, upper (is the upper half of a horizontal) - make horizontal, upper (will be upper half of a vertical)
				b.parent = &subSpace[A]{
					center:   b.center.Minus(geo.V(b.radius.X, 0)),
					radius:   b.radius.Times(geo.V(2, 1)),
					vertical: false,
					upper:    true,
					children: [2]*subSpace[A]{nil, b},
					size:     b.size,
				}
			}
		}
		b.parent.clean()
	}
	return b.parent
}

func (b *subSpace[A]) demandChild(upper bool) *subSpace[A] {
	if !upper {
		if b.children[0] == nil {
			if !b.vertical {
				radius := b.radius.Over(geo.V(2, 1))
				b.children[0] = &subSpace[A]{
					parent:   b,
					center:   b.center.Minus(geo.V(radius.X, 0)),
					radius:   radius,
					vertical: true,
					upper:    false,
				}
			} else {
				radius := b.radius.Over(geo.V(1, 2))
				b.children[0] = &subSpace[A]{
					parent:   b,
					center:   b.center.Minus(geo.V(0, radius.Y)),
					radius:   radius,
					vertical: false,
					upper:    false,
				}
			}
		}
		return b.children[0]
	}
	if b.children[1] == nil {
		if !b.vertical {
			radius := b.radius.Over(geo.V(2, 1))
			b.children[1] = &subSpace[A]{
				parent:   b,
				center:   b.center.Plus(geo.V(radius.X, 0)),
				radius:   radius,
				vertical: true,
				upper:    true,
			}
		} else {
			radius := b.radius.Over(geo.V(1, 2))
			b.children[1] = &subSpace[A]{
				parent:   b,
				center:   b.center.Plus(geo.V(0, radius.Y)),
				radius:   radius,
				vertical: false,
				upper:    true,
			}
		}
	}
	return b.children[1]
}
