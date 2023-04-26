package geometry

import (
	"math"
)

type Vector struct {
	X, Y float64
}

func V(x, y float64) Vector { return Vector{x, y} }

func (v Vector) Plus(v2 Vector) Vector {
	return Vector{v.X + v2.X, v.Y + v2.Y}
}

func (v Vector) Minus(v2 Vector) Vector {
	return Vector{v.X - v2.X, v.Y - v2.Y}
}

func (v Vector) Times(v2 Vector) Vector {
	return Vector{v.X * v2.X, v.Y * v2.Y}
}

func (v Vector) Over(v2 Vector) Vector {
	return Vector{v.X / v2.X, v.Y / v2.Y}
}

func (v Vector) PlusScalar(f float64) Vector {
	return Vector{v.X + f, v.Y + f}
}

func (v Vector) MinusScalar(f float64) Vector {
	return Vector{v.X - f, v.Y - f}
}

func (v Vector) TimesScalar(f float64) Vector {
	return Vector{v.X * f, v.Y * f}
}

func (v Vector) OverScalar(f float64) Vector {
	return Vector{v.X / f, v.Y / f}
}

func (v Vector) Dot(v2 Vector) float64 {
	return v.X*v2.X + v.Y*v2.Y
}

func (v Vector) Cross(v2 Vector) float64 {
	return v.X*v2.Y - v.Y*v2.X
}

func (v Vector) MagnitudeSquared() float64 { return v.Dot(v) }

func (v Vector) Magnitude() float64 { return math.Sqrt(v.MagnitudeSquared()) }

func (v Vector) Floored() Vector { return Vector{math.Floor(v.X), math.Floor(v.Y)} }

func (v Vector) ShapeType() ShapeType { return POINT }

func (v Vector) Vertex(i int) Vector {
	if i != 0 {
		panic("index out of bounds")
	}
	return v
}

func (v Vector) Bounds() Rectangle { return Rectangle{v, v} }
