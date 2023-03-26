package space

import (
	"fmt"
	"math"

	geo "github.com/clayts/gec/geometry"
)

const (
	maxCached = 1
)

type Space[T any] struct {
	local     *subSpace[T]
	universal *subSpace[T]
}

func New[T any]() *Space[T] { return &Space[T]{} }

// Runs f(*Leaf) on every *Leaf in the Tree which intersects with s.
func (spc *Space[T]) AllIntersecting(s geo.Shape, f func(z *Zone[T]) bool) bool {
	return spc.universal.all(f) && spc.local.allIntersecting(s, f)
}

// Runs f(*Leaf) on every *Leaf in the Tree.
func (spc *Space[T]) All(f func(z *Zone[T]) bool) bool {
	return spc.universal.all(f) && spc.local.all(f)
}

func (spc *Space[T]) add(z *Zone[T]) {
	if z.shape == geo.EVERYWHERE {
		if spc.universal == nil {
			spc.universal = &subSpace[T]{}
		}
		spc.universal.size++
		spc.universal.store(z)
	} else {
		if spc.local == nil {
			spc.local = &subSpace[T]{radius: geo.V(math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64)}
		}
		for !spc.local.bounds().Contains(z.shape.Bounds()) {
			spc.local = spc.local.demandParent()
		}
		spc.local.add(z)
	}
}

func (spc *Space[T]) remove(z *Zone[T]) {
	z.subSpace.remove(z)

	if spc.local != nil && spc.local.size == 0 {
		spc.local = nil
	}
	if spc.universal != nil && spc.universal.size == 0 {
		spc.universal = nil
	}
}

func (spc *Space[T]) Size() int {
	size := 0
	if spc.local != nil {
		size += spc.local.size
	}
	if spc.universal != nil {
		size += spc.universal.size
	}
	return size
}

func (spc *Space[T]) PrintStructure() {
	var print func(b *subSpace[T], prefix string)
	print = func(b *subSpace[T], prefix string) {
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
