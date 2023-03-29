package geometry

import (
	"math"
)

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

func AllEdges(s Shape, f func(g Segment) bool) bool {
	if s.ShapeType() == POINT {
		return true
	}
	a := s.Vertex(Len(s) - 1)
	return AllVertices(s, func(i int, b Vector) bool {
		if !f(Segment{a, b}) {
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

// AllClockwise returns true if all of s is clockwise of g[1], from the perspective of g[0]
func AllClockwise(g Segment, s Shape) bool {
	line := g[1].Minus(g[0])
	return AllVertices(s, func(i int, v Vector) bool { return v.Minus(g[0]).Cross(line) > 0 })
}

// AllAnticlockwise returns true if all of s is anticlockwise of g[1], from the perspective of g[0]
func AllAnticlockwise(g Segment, s Shape) bool {
	line := g[1].Minus(g[0])
	return AllVertices(s, func(i int, v Vector) bool { return v.Minus(g[0]).Cross(line) < 0 })
}

func Contains(s, s2 Shape) bool {
	return s.Bounds().Contains(s2.Bounds()) && ContainsSkipBoundsCheck(s, s2)

}

func ContainsSkipBoundsCheck(s, s2 Shape) bool {
	return s.ShapeType() == RECTANGLE || AllEdges(s, func(g Segment) bool { return g.AxisAligned() || AllClockwise(g, s2) })
}

func Intersects(s, s2 Shape) bool {
	return s.Bounds().Intersects(s2.Bounds()) && IntersectsSkipBoundsCheck(s, s2)
}

func IntersectsSkipBoundsCheck(s, s2 Shape) bool {
	return (s.ShapeType() == RECTANGLE || AllEdges(s, func(g Segment) bool { return g.AxisAligned() || !AllAnticlockwise(g, s2) })) &&
		(s2.ShapeType() == RECTANGLE || AllEdges(s2, func(g Segment) bool { return g.AxisAligned() || !AllAnticlockwise(g, s) }))
}
