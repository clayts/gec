package geometry

type Polygon struct {
	vertices []Vector
	bounds   Rectangle
}

func P(vertices ...Vector) Polygon {
	if len(vertices) < 3 {
		panic("polygons must have at least 3 vertices")
	}

	p := Polygon{
		vertices: vertices,
	}

	for i, v := range p.vertices {
		if i == 0 {
			p.bounds.Min = v
			p.bounds.Max = v
		} else {
			if v.X < p.bounds.Min.X {
				p.bounds.Min.X = v.X
			} else if v.X > p.bounds.Max.X {
				p.bounds.Max.X = v.X
			}
			if v.Y < p.bounds.Min.Y {
				p.bounds.Min.Y = v.Y
			} else if v.Y > p.bounds.Max.Y {
				p.bounds.Max.Y = v.Y
			}
		}
	}

	return p
}

func (p Polygon) Vertex(i int) Vector { return p.vertices[i] }

func (p Polygon) ShapeType() ShapeType { return ShapeType(len(p.vertices)) }

func (p Polygon) Bounds() Rectangle { return p.bounds }
