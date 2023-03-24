package tree

import (
	"fmt"
	"math"

	geo "github.com/clayts/gec/geometry"
)

const (
	maxCached = 1
)

type Tree struct {
	local     *branch
	universal *branch
}

// Runs f(*Leaf) on every *Leaf in the Tree which intersects with s.
func (t *Tree) AllLeavesIntersecting(s geo.Shape, f func(l *Leaf) bool) bool {
	return t.universal.allLeaves(f) && t.local.allLeavesIntersecting(s, f)
}

// Runs f(*Leaf) on every *Leaf in the Tree.
func (t *Tree) AllLeaves(f func(l *Leaf) bool) bool {
	return t.universal.allLeaves(f) && t.local.allLeaves(f)
}

func (t *Tree) add(l *Leaf) {
	if l.shape == geo.EVERYWHERE {
		if t.universal == nil {
			t.universal = &branch{}
		}
		t.universal.size++
		t.universal.store(l)
	} else {
		if t.local == nil {
			t.local = &branch{radius: geo.V(math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64)}
		}
		for !t.local.bounds().Contains(l.shape.Bounds()) {
			t.local = t.local.demandParent()
		}
		t.local.add(l)
	}
}

func (t *Tree) remove(l *Leaf) {
	l.branch.remove(l)

	if t.local != nil && t.local.size == 0 {
		t.local = nil
	}
	if t.universal != nil && t.universal.size == 0 {
		t.universal = nil
	}
}

// DEBUG -----------------------------------------------------------
func (t *Tree) Print() {
	var print func(b *branch, prefix string)
	print = func(b *branch, prefix string) {
		if b != nil {
			if b.upper {
				prefix += "1"
			} else {
				prefix += "0"
			}
			fmt.Println(prefix, b.size, "(", len(b.leaves), ")")
			for _, c := range b.children {
				print(c, prefix)
			}
		}
	}
	print(t.local, "L")
	print(t.universal, "U")
}

func (t *Tree) Size() int {
	size := 0
	if t.local != nil {
		size += t.local.size
	}
	if t.universal != nil {
		size += t.universal.size
	}
	return size
}

// -----------------------------------------------------------------
