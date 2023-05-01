package graphics

type Layout []struct {
	Location                     AttributeLocation
	Elements, Components, Length int32
}

type autoVBO struct {
	vbo          VBO
	vboCapacity  int
	stride       int32
	data         []float32
	synchronized bool
}

func openAutoVBO() *autoVBO {
	b := &autoVBO{}
	b.vbo = OpenVBO()

	return b
}

func (avb *autoVBO) attach(vao VAO, divisor uint32, layout Layout) {
	info := []struct {
		location AttributeLocation
		offset   int32
		size     int32
	}{}

	for _, attribute := range layout {
		location := attribute.Location
		for i := int32(0); i < attribute.Elements; i++ {
			for j := int32(0); j < attribute.Components; j++ {
				info = append(info, struct {
					location AttributeLocation
					offset   int32
					size     int32
				}{location, avb.stride, attribute.Length})
				location++
				avb.stride += attribute.Length
			}
		}
	}

	for _, a := range info {
		vao.SetAttribute(a.location, a.offset, a.size, avb.stride, divisor, avb.vbo)
	}
}

func (avb *autoVBO) Add(data ...float32) {
	avb.data = append(avb.data, data...)
	avb.synchronized = false
}

func (avb *autoVBO) Clear() {
	avb.data = avb.data[:0]
	avb.synchronized = false
}

func (avb *autoVBO) Length() int32 {
	if avb == nil {
		return 0
	}
	return int32(len(avb.data))
}

func (avb *autoVBO) Stride() int32 { return avb.stride }

func (avb *autoVBO) sync() {
	if !avb.synchronized {
		if len(avb.data) <= avb.vboCapacity {
			avb.vbo.SetDataSubsection(0, avb.data...)
		} else {
			avb.vbo.SetData(avb.data...)
			avb.vboCapacity = len(avb.data)
		}
		avb.synchronized = true
	}
}

func (avb *autoVBO) close() {
	avb.vbo.Close()
}

type Buffer struct {
	mode                Mode
	vao                 VAO
	vertices, instances *autoVBO
}

func OpenBuffer(mode Mode, vertexLayout, insanceLayout Layout) Buffer {
	r := Buffer{}
	r.mode = mode
	r.vao = OpenVAO()
	r.vertices = openAutoVBO()
	r.vertices.attach(r.vao, 0, vertexLayout)
	if insanceLayout != nil {
		r.instances = openAutoVBO()
		r.instances.attach(r.vao, 1, insanceLayout)
	}
	return r
}

func (buf Buffer) Vertices() interface {
	Add(data ...float32)
	Clear()
	Length() int32
	Stride() int32
} {
	return buf.vertices
}

func (buf Buffer) Instances() interface {
	Add(data ...float32)
	Clear()
	Length() int32
	Stride() int32
} {
	return buf.instances
}

func (buf Buffer) Draw(program Program) {
	if vlen := buf.vertices.Length(); vlen > 0 {
		buf.vertices.sync()
		if ilen := buf.instances.Length(); ilen > 0 {
			buf.instances.sync()
			buf.vao.DrawInstanced(program, buf.mode, 0, vlen/buf.vertices.stride, ilen/buf.instances.stride)
		} else {
			buf.vao.Draw(program, buf.mode, 0, vlen/buf.vertices.stride)
		}
	}
}

func (buf Buffer) Close() {
	buf.vao.Close()
	buf.vertices.close()
	buf.instances.close()
}
