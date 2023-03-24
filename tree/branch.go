package tree

import (
	geo "github.com/clayts/gec/geometry"
)

type branch struct {
	parent         *branch
	center, radius geo.Vector
	children       [2]*branch
	vertical       bool
	upper          bool
	leaves         []*Leaf
	size           int
}

func (b *branch) allContained(s geo.Shape, f func(l *Leaf) bool) bool {
	if b == nil {
		return true
	}
	for _, l := range b.leaves {
		if (s.ShapeType() == geo.RECTANGLE || geo.AllEdges(s, func(i int, g geo.Segment) bool { return g.AxisAligned() || !geo.LeftOf(g, l.shape) })) &&
			(l.shape.ShapeType() == geo.RECTANGLE || geo.AllEdges(l.shape, func(i int, g geo.Segment) bool { return g.AxisAligned() || !geo.LeftOf(g, s) })) {
			if !f(l) {
				return false
			}
		}
	}
	return b.children[0].allContained(s, f) && b.children[1].allContained(s, f)
}

func (b *branch) all(s geo.Shape, f func(l *Leaf) bool) bool {
	if b == nil {
		return true
	}
	if !s.Bounds().Intersects(b.bounds()) {
		return true
	}
	if s.Bounds().Contains(b.bounds()) {
		return b.allContained(s, f)
	}
	for _, l := range b.leaves {
		if geo.Intersects(s, l.shape) {
			if !f(l) {
				return false
			}
		}
	}
	return b.children[0].all(s, f) && b.children[1].all(s, f)
}

func (b *branch) remove(l *Leaf) {
	// delete from list
	finalIndex := len(b.leaves) - 1
	b.leaves[l.index] = b.leaves[finalIndex]
	b.leaves[l.index].index = l.index
	b.leaves[finalIndex] = nil
	b.leaves = b.leaves[:finalIndex]

	// update leaf
	l.branch = nil

	b.decrease()
}

func (b *branch) clean() {
	if b.size <= maxCached {
		if b.children[0] != nil {
			for _, l := range b.children[0].leaves {
				b.store(l)
			}
			b.children[0] = nil
		}
		if b.children[1] != nil {
			for _, l := range b.children[1].leaves {
				b.store(l)
			}
			b.children[1] = nil
		}
	}
}

func (b *branch) decrease() {
	b.size--
	b.clean()
	if b.parent != nil {
		b.parent.decrease()
	}
}

func (b *branch) add(l *Leaf) {
	// if size < max, cache this
	// if size == max, distribute all
	// if size > max,  distribute this
	if b.size < maxCached {
		b.store(l)
	} else {
		if b.size == maxCached {
			leaves := b.leaves
			b.leaves = b.leaves[:0]
			for _, l2 := range leaves {
				b.distribute(l2)
			}
		}
		b.distribute(l)
	}
	b.size++
}

func (b *branch) distribute(l *Leaf) {
	if !b.vertical {
		if l.shape.Bounds().Max.X < b.center.X {
			b.demandChild(false).add(l)
			return
		} else if l.shape.Bounds().Min.X > b.center.X {
			b.demandChild(true).add(l)
			return
		}
	} else {
		if l.shape.Bounds().Max.Y < b.center.Y {
			b.demandChild(false).add(l)
			return
		} else if l.shape.Bounds().Min.Y > b.center.Y {
			b.demandChild(true).add(l)
			return
		}
	}
	b.store(l)
}

func (b *branch) store(l *Leaf) {
	l.index = len(b.leaves)
	l.branch = b
	b.leaves = append(b.leaves, l)
}

func (b *branch) bounds() geo.Rectangle {
	return geo.R(b.center.Minus(b.radius), b.center.Plus(b.radius))
}

func (b *branch) demandParent() *branch {
	if b.parent == nil {
		if !b.vertical {
			if !b.upper {
				// horizontal, lower (is the lower half of a vertical) - make vertical, upper (will be upper half of a horizontal)
				b.parent = &branch{
					center:   b.center.Plus(geo.V(0, b.radius.Y)),
					radius:   b.radius.Times(geo.V(1, 2)),
					vertical: true,
					upper:    true,
					children: [2]*branch{b, nil},
					size:     b.size,
				}
			} else {
				// horizontal, upper (is the upper half of a vertical) - make vertical, lower (will be lower half of a horizontal)
				b.parent = &branch{
					center:   b.center.Minus(geo.V(0, b.radius.Y)),
					radius:   b.radius.Times(geo.V(1, 2)),
					vertical: true,
					upper:    false,
					children: [2]*branch{nil, b},
					size:     b.size,
				}
			}
		} else {
			if !b.upper {
				// vertical, lower  (is the lower half of a horizontal) - make horizontal, lower (will be lower half of a vertical)
				b.parent = &branch{
					center:   b.center.Plus(geo.V(b.radius.X, 0)),
					radius:   b.radius.Times(geo.V(2, 1)),
					vertical: false,
					upper:    false,
					children: [2]*branch{b, nil},
					size:     b.size,
				}
			} else {
				// vertical, upper (is the upper half of a horizontal) - make horizontal, upper (will be upper half of a vertical)
				b.parent = &branch{
					center:   b.center.Minus(geo.V(b.radius.X, 0)),
					radius:   b.radius.Times(geo.V(2, 1)),
					vertical: false,
					upper:    true,
					children: [2]*branch{nil, b},
					size:     b.size,
				}
			}
		}
		b.parent.clean()
	}
	return b.parent
}

func (b *branch) demandChild(upper bool) *branch {
	if !upper {
		if b.children[0] == nil {
			if !b.vertical {
				radius := b.radius.Over(geo.V(2, 1))
				b.children[0] = &branch{
					parent:   b,
					center:   b.center.Minus(geo.V(radius.X, 0)),
					radius:   radius,
					vertical: true,
					upper:    false,
				}
			} else {
				radius := b.radius.Over(geo.V(1, 2))
				b.children[0] = &branch{
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
			b.children[1] = &branch{
				parent:   b,
				center:   b.center.Plus(geo.V(radius.X, 0)),
				radius:   radius,
				vertical: true,
				upper:    true,
			}
		} else {
			radius := b.radius.Over(geo.V(1, 2))
			b.children[1] = &branch{
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
