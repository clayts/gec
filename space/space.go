package space

import (
	"fmt"
	"math"

	geo "github.com/clayts/gec/geometry"
)

const (
	maxCached = 1
)

type Space[A any] struct {
	local     *subSpace[A]
	universal *subSpace[A]
}

func New[A any]() *Space[A] { return &Space[A]{} }

// Runs f(*Leaf) on every *Leaf in the Tree which intersects with s.
func (spc *Space[A]) AllIntersecting(s geo.Shape, f func(z *Zone[A]) bool) bool {
	return spc.universal.all(f) && spc.local.allIntersecting(s, f)
}

// Runs f(*Leaf) on every *Leaf in the Tree.
func (spc *Space[A]) All(f func(z *Zone[A]) bool) bool {
	return spc.universal.all(f) && spc.local.all(f)
}

func (spc *Space[A]) add(z *Zone[A]) {
	if z.shape == geo.EVERYWHERE {
		if spc.universal == nil {
			spc.universal = &subSpace[A]{}
		}
		spc.universal.size++
		spc.universal.store(z)
	} else {
		if spc.local == nil {
			spc.local = &subSpace[A]{radius: geo.V(math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64)}
		}
		for !spc.local.bounds().Contains(z.shape.Bounds()) {
			spc.local = spc.local.demandParent()
		}
		spc.local.add(z)
	}
}

func (spc *Space[A]) remove(z *Zone[A]) {
	z.subSpace.remove(z)

	if spc.local != nil && spc.local.size == 0 {
		spc.local = nil
	}
	if spc.universal != nil && spc.universal.size == 0 {
		spc.universal = nil
	}
}

func (spc *Space[A]) Size() int {
	size := 0
	if spc.local != nil {
		size += spc.local.size
	}
	if spc.universal != nil {
		size += spc.universal.size
	}
	return size
}

func (spc *Space[A]) PrintStructure() {
	var print func(b *subSpace[A], prefix string)
	print = func(b *subSpace[A], prefix string) {
		if b != nil {
			if b.upper {
				prefix += "1"
			} else {
				prefix += "0"
			}
			fmt.Println(prefix, b.size, "(", len(b.zones), ")")
			for _, c := range b.children {
				print(c, prefix)
			}
		}
	}
	print(spc.local, "L")
	print(spc.universal, "U")
}
