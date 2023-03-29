package set

func (s *Set[T]) New() *Entity[T] {
	e := &Entity[T]{set: s, index: -1}
	return e
}

func (e *Entity[T]) Enable() *Entity[T] {
	e.set.add(e)
	return e
}

func (e *Entity[T]) Disable() *Entity[T] {
	e.set.remove(e)
	return e
}

func (e *Entity[T]) SetContents(c T) *Entity[T] {
	e.Contents = c
	return e
}
