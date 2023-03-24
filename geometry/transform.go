package geometry

type Transform [2][3]float64

func T() Transform {
	return Transform{
		{1, 0, 0},
		{0, 1, 0},
	}
}

func Translation(translation Vector) Transform {
	return Transform{
		{1, 0, translation.X},
		{0, 1, translation.Y},
	}
}

func Scale(scale Vector) Transform {
	return Transform{
		{scale.X, 0, 0},
		{0, scale.Y, 0},
	}
}

func Rotation(rotation Angle) Transform {
	s, c := rotation.SinCos()
	return Transform{
		{c, -s, 0},
		{s, c, 0},
	}
}

// ts[n-1].Times(ts[n-2].Times(...))
func TransformSequence(ts ...Transform) Transform {
	if len(ts) == 0 {
		return T()
	}
	if len(ts) == 1 {
		return ts[0]
	}
	t := ts[0]
	for i := 1; i < len(ts); i++ {
		t = ts[i].Times(t)
	}
	return t
}

func RotationAround(rotation Angle, point Vector) Transform {
	return TransformSequence(
		Translation(Vector{0, 0}.Minus(point)),
		Rotation(rotation),
		Translation(point),
	)
}

func Reframe(from, to Rectangle) Transform {
	scale := to.Size().Over(from.Size())
	translation := to.Min.Minus(from.Min)
	return Transform{
		{scale.X, 0, translation.X},
		{0, scale.Y, translation.Y},
	}
}

func (t Transform) Times(t2 Transform) Transform {
	return Transform{
		{t[0][0]*t2[0][0] + t[0][1]*t2[1][0], t[0][0]*t2[0][1] + t[0][1]*t2[1][1], t[0][0]*t2[0][2] + t[0][1]*t2[1][2] + t[0][2]},
		{t[1][0]*t2[0][0] + t[1][1]*t2[1][0], t[1][0]*t2[0][1] + t[1][1]*t2[1][1], t[1][0]*t2[0][2] + t[1][1]*t2[1][2] + t[1][2]},
	}
}

func (t Transform) Inverse() Transform {
	det := 1 / (t[0][0]*t[1][1] - t[1][0]*t[0][1])
	return Transform{
		{t[1][1] * det, -t[0][1] * det, (t[0][1]*t[1][2] - t[1][1]*t[0][2]) * det},
		{-t[1][0] * det, t[0][0] * det, (t[1][0]*t[0][2] - t[0][0]*t[1][2]) * det},
	}
}

func (t Transform) Vector(v Vector) Vector {
	return Vector{
		t[0][0]*v.X + t[0][1]*v.Y + t[0][2],
		t[1][0]*v.X + t[1][1]*v.Y + t[1][2],
	}
}

func (t Transform) Segment(g Segment) Segment {
	return Segment{t.Vector(g[0]), t.Vector(g[1])}
}

func (t Transform) Rectangle(r Rectangle) Shape {
	if t[0][1] == 0 && t[1][0] == 0 {
		return Rectangle{t.Vector(r.Min), t.Vector(r.Max)}
	}
	return t.shapeToPolygon(r)
}

func (t Transform) Polygon(p Polygon) Polygon {
	return t.shapeToPolygon(p)
}

func (t Transform) Shape(s Shape) Shape {
	if s.ShapeType() == RECTANGLE {
		return t.Rectangle(s.Bounds())
	}
	if s.ShapeType() == SEGMENT {
		return t.Segment(Segment{s.Vertex(0), s.Vertex(1)})
	}
	if s.ShapeType() == POINT {
		return t.Vector(s.Vertex(0))
	}
	return t.shapeToPolygon(s)
}

func (t Transform) shapeToPolygon(s Shape) Polygon {
	vs := make([]Vector, Len(s))
	AllVertices(s, func(i int, v Vector) bool {
		vs[i] = t.Vector(v)
		return true
	})
	return P(vs...)
}
