package geometry

type Rectangle struct {
	Min, Max Vector
}

func R(min, max Vector) Rectangle { return Rectangle{Min: min, Max: max} }

func (r Rectangle) Size() Vector { return r.Max.Minus(r.Min) }

func (r Rectangle) Radius() Vector { return r.Size().Over(Vector{2, 2}) }

func (r Rectangle) Center() Vector { return r.Max.Plus(r.Min).Over(Vector{2, 2}) }

func (r Rectangle) ShapeType() ShapeType { return RECTANGLE }

func (r Rectangle) Vertex(i int) Vector {
	switch i {
	case 0:
		return r.Min
	case 1:
		return Vector{X: r.Min.X, Y: r.Max.Y}
	case 2:
		return r.Max
	case 3:
		return Vector{X: r.Max.X, Y: r.Min.Y}
	}
	panic("index out of bounds")
}

func (r Rectangle) Bounds() Rectangle { return r }

func (r Rectangle) Contains(r2 Rectangle) bool {
	return r.Min.X < r2.Min.X &&
		r.Max.X > r2.Max.X &&
		r.Min.Y < r2.Min.Y &&
		r.Max.Y > r2.Max.Y
}

func (r Rectangle) Intersects(r2 Rectangle) bool {
	return r.Max.X > r2.Min.X &&
		r.Min.X < r2.Max.X &&
		r.Max.Y > r2.Min.Y &&
		r.Min.Y < r2.Max.Y
}
