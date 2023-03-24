package geometry

import "math"

const (
	RECTANGLE ShapeType = iota
	POINT
	SEGMENT
)

var (
	EVERYWHERE Rectangle = Rectangle{Vector{math.Inf(-1), math.Inf(-1)}, Vector{math.Inf(1), math.Inf(1)}}
)

type ShapeType int

type Shape interface {
	Vertex(i int) Vector
	ShapeType() ShapeType
	Bounds() Rectangle
}

func AllVertices(s Shape, f func(i int, v Vector) bool) bool {
	for i := 0; i < Len(s); i++ {
		if !f(i, s.Vertex(i)) {
			return false
		}
	}
	return true
}

func AllEdges(s Shape, f func(i int, g Segment) bool) bool {
	if s.ShapeType() == POINT {
		return true
	}
	a := s.Vertex(Len(s) - 1)
	return AllVertices(s, func(i int, b Vector) bool {
		if !f(i, Segment{a, b}) {
			return false
		}
		a = b
		return true
	})
}

func Len(s Shape) int {
	if s.ShapeType() == RECTANGLE {
		return 4
	}
	return int(s.ShapeType())
}

// RightOf returns true if all of s is to the left of g[1], from the perspective of g[0]
func RightOf(g Segment, s Shape) bool {
	line := g[1].Minus(g[0])
	return AllVertices(s, func(i int, v Vector) bool { return v.Minus(g[0]).Cross(line) <= 0 })
}

// LeftOf returns true if all of s is to the right of g[1], from the perspective of g[0]
func LeftOf(g Segment, s Shape) bool {
	line := g[1].Minus(g[0])
	return AllVertices(s, func(i int, v Vector) bool { return v.Minus(g[0]).Cross(line) >= 0 })
}

func Contains(s, s2 Shape) bool {
	return s.Bounds().Contains(s2.Bounds()) &&
		(s.ShapeType() == RECTANGLE || AllEdges(s, func(i int, g Segment) bool { return g.AxisAligned() || RightOf(g, s2) }))
}

func Intersects(s, s2 Shape) bool {
	return s.Bounds().Intersects(s2.Bounds()) &&
		(s.ShapeType() == RECTANGLE || AllEdges(s, func(i int, g Segment) bool { return g.AxisAligned() || !LeftOf(g, s2) })) &&
		(s2.ShapeType() == RECTANGLE || AllEdges(s2, func(i int, g Segment) bool { return g.AxisAligned() || !LeftOf(g, s) }))
}
