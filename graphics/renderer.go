package graphics

type buffer struct {
	vbo         VBO
	vboCapacity int
	vboLength   int
	stride      int32
}

func openBuffer(vao VAO, attributes []Attribute, divisor uint32, layout ...string) *buffer {
	b := &buffer{}
	b.vbo = OpenVBO()

	layoutData := []struct {
		location AttributeLocation
		offset   int32
		size     int32
	}{}

	for _, name := range layout {
		var attribute Attribute
		for _, a := range attributes {
			if a.Name() == name {
				attribute = a
				break
			}
		}
		if attribute.Name() == "" {
			panic("no such attribute")
		}
		location := attribute.Location()
		l, m, n := attribute.Size()
		for i := int32(0); i < l; i++ {
			// iterate array
			for j := int32(0); j < m; j++ {
				// iterate m
				layoutData = append(layoutData, struct {
					location AttributeLocation
					offset   int32
					size     int32
				}{location, b.stride, n})
				location++
				b.stride += n
			}
		}
	}

	for _, a := range layoutData {
		vao.SetAttribute(a.location, a.offset, a.size, b.stride, divisor, b.vbo)
	}
	return b

}

func (b *buffer) set(data ...float32) {
	b.vboLength = len(data)
	if b.vboLength <= b.vboCapacity {
		b.vbo.SetDataSubsection(0, data...)
	} else {
		b.vbo.SetData(data...)
		b.vboCapacity = b.vboLength
	}
}

func (b *buffer) count() int32 {
	return int32(b.vboLength) / b.stride
}

func (b *buffer) close() {
	b.vbo.Close()
}

type Renderer struct {
	mode      Mode
	program   Program
	vao       VAO
	vertices  *buffer
	instances *buffer
}

func OpenRenderer(mode Mode, vertexLayout []string, instanceLayout []string, program Program) *Renderer {
	r := &Renderer{}
	r.mode = mode
	r.program = program
	r.vao = NewVAO()
	if len(vertexLayout) > 0 {
		r.vertices = openBuffer(r.vao, r.program.Attributes(), 0, vertexLayout...)
	}
	if len(instanceLayout) > 0 {
		r.instances = openBuffer(r.vao, r.program.Attributes(), 1, instanceLayout...)
	}
	return r
}

func (r *Renderer) SetVertices(vertices ...float32) {
	r.vertices.set(vertices...)
}

func (r *Renderer) SetInstances(instances ...float32) {
	r.instances.set(instances...)
}

func (r *Renderer) Render() {
	if r.vertices != nil {
		if r.instances != nil {
			if vc, ic := r.vertices.count(), r.instances.count(); vc > 0 && ic > 0 {
				r.program.DrawInstanced(r.vao, r.mode, 0, vc, ic)
			}
		} else {
			if vc := r.vertices.count(); vc > 0 {
				r.program.Draw(r.vao, r.mode, 0, vc)
			}
		}
	}
}

func (r *Renderer) Close() {
	r.vao.Close()
	r.vertices.close()
	r.instances.close()
}
