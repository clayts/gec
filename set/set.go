package set

type Set[A any] struct {
	entities []*Entity[A]
}

func New[A any]() *Set[A] {
	s := &Set[A]{}
	return s
}
func (s *Set[A]) All(f func(e *Entity[A]) bool) bool {
	for i := len(s.entities) - 1; i >= 0; i-- {
		if !f(s.entities[i]) {
			return false
		}
	}
	return true
}

func (s *Set[A]) Size() int { return len(s.entities) }

func (s *Set[A]) add(e *Entity[A]) {
	if e.index == -1 {
		e.index = len(s.entities)
		s.entities = append(s.entities, e)
	}
}

func (s *Set[A]) remove(e *Entity[A]) {
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
