package set

type Entity[A any] struct {
	set      *Set[A]
	index    int
	Contents A
}

func (s *Set[A]) New() *Entity[A] {
	e := &Entity[A]{set: s, index: -1}
	return e
}

func (e *Entity[A]) SetState(state bool) *Entity[A] {
	if state {
		e.set.add(e)
	} else {
		e.set.remove(e)
	}
	return e
}

func (e *Entity[A]) SetContents(c A) *Entity[A] {
	e.Contents = c
	return e
}
