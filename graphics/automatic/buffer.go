package automatic

import "github.com/clayts/gec/graphics"

type Layout []struct {
	Location                     graphics.AttributeLocation
	Elements, Components, Length int32
}

type subBuffer struct {
	vbo          graphics.VBO
	vboCapacity  int
	stride       int32
	data         []float32
	synchronized bool
}

func openSubBuffer() *subBuffer {
	b := &subBuffer{}
	b.vbo = graphics.OpenVBO()

	return b
}

func (sub *subBuffer) attach(vao graphics.VAO, divisor uint32, layout Layout) {
	info := []struct {
		location graphics.AttributeLocation
		offset   int32
		size     int32
	}{}

	for _, attribute := range layout {
		location := attribute.Location
		for i := int32(0); i < attribute.Elements; i++ {
			for j := int32(0); j < attribute.Components; j++ {
				info = append(info, struct {
					location graphics.AttributeLocation
					offset   int32
					size     int32
				}{location, sub.stride, attribute.Length})
				location++
				sub.stride += attribute.Length
			}
		}
	}

	for _, a := range info {
		vao.SetAttribute(a.location, a.offset, a.size, sub.stride, divisor, sub.vbo)
	}
}

func (sub *subBuffer) add(data ...float32) {
	sub.data = append(sub.data, data...)
	sub.synchronized = false
}

func (sub *subBuffer) clear() {
	sub.data = sub.data[:0]
	sub.synchronized = false
}

func (sub *subBuffer) length() int32 {
	if sub == nil {
		return 0
	}
	return int32(len(sub.data))
}

func (sub *subBuffer) sync() {
	if !sub.synchronized {
		if len(sub.data) <= sub.vboCapacity {
			sub.vbo.SetDataSubsection(0, sub.data...)
		} else {
			sub.vbo.SetData(sub.data...)
			sub.vboCapacity = len(sub.data)
		}
		sub.synchronized = true
	}
}

func (sub *subBuffer) close() {
	sub.vbo.Close()
}

type Buffer struct {
	mode                graphics.Mode
	vao                 graphics.VAO
	vertices, instances *subBuffer
}

func OpenBuffer(mode graphics.Mode, vertexLayout, insanceLayout Layout) Buffer {
	r := Buffer{}
	r.mode = mode
	r.vao = graphics.OpenVAO()
	r.vertices = openSubBuffer()
	r.vertices.attach(r.vao, 0, vertexLayout)
	if insanceLayout != nil {
		r.instances = openSubBuffer()
		r.instances.attach(r.vao, 1, insanceLayout)
	}
	return r
}

func (buf Buffer) Draw(program graphics.Program) {
	if vlen := buf.vertices.length(); vlen > 0 {
		buf.vertices.sync()
		if ilen := buf.instances.length(); ilen > 0 {
			buf.instances.sync()
			buf.vao.DrawInstanced(program, buf.mode, 0, vlen/buf.vertices.stride, ilen/buf.instances.stride)
		} else {
			buf.vao.Draw(program, buf.mode, 0, vlen/buf.vertices.stride)
		}
	}
}

func (buf Buffer) Clear(vertices, instances bool) {
	if vertices {
		buf.vertices.clear()
	}
	if instances {
		buf.instances.clear()
	}
}

func (buf Buffer) Close() {
	buf.vao.Close()
	buf.vertices.close()
	buf.instances.close()
}

type Vertex []float32

func (vtx Vertex) Draw(buf Buffer) {
	buf.vertices.add(vtx...)
}

type Instance []float32

func (ins Instance) Draw(buf Buffer) {
	buf.instances.add(ins...)
}
