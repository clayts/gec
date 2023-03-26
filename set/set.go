package set

type Entity[T any] struct {
	set      *Set[T]
	index    int
	Contents T
}

type Set[T any] struct {
	entities []*Entity[T]
}

func New[T any]() *Set[T] {
	s := &Set[T]{}
	return s
}
func (s *Set[T]) All(f func(e *Entity[T]) bool) bool {
	for i := len(s.entities) - 1; i >= 0; i-- {
		if !f(s.entities[i]) {
			return false
		}
	}
	return true
}

func (s *Set[T]) Size() int { return len(s.entities) }

func (s *Set[T]) add(e *Entity[T]) {
	if e.index == -1 {
		e.index = len(s.entities)
		s.entities = append(s.entities, e)
	}
}

func (s *Set[T]) remove(e *Entity[T]) {
	if e.index != -1 {
		finalIndex := len(s.entities) - 1
		if finalIndex != 0 {
			finalEntity := s.entities[finalIndex]
			s.entities[e.index] = finalEntity
			finalEntity.index = e.index
			s.entities[finalIndex] = nil
		}
		s.entities = s.entities[:finalIndex]
		e.index = -1
	}
}
