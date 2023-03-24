package geometry

type Segment [2]Vector

func S(a, b Vector) Segment { return Segment{a, b} }

func (g Segment) AxisAligned() bool {
	return g[0].X == g[1].X || g[0].Y == g[1].Y
}

func (g Segment) ShapeType() ShapeType { return SEGMENT }

func (g Segment) Vertex(i int) Vector { return g[i] }

func (g Segment) Bounds() Rectangle {
	r := Rectangle{g[0], g[1]}
	if r.Min.X > r.Max.X {
		r.Min.X, r.Max.X = r.Max.X, r.Min.X
	}
	if r.Min.Y > r.Max.Y {
		r.Min.Y, r.Max.Y = r.Max.Y, r.Min.Y
	}
	return r
}
