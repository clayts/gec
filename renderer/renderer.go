package renderer

import (
	"strconv"

	gfx "github.com/clayts/gec/graphics"
)

type buffer struct {
	vbo     gfx.VBO
	vboSize int
	stride  int32
	data    []float32
	changed bool
}

func newBuffer(vao gfx.VAO, attributes []gfx.Attribute, layout ...string) *buffer {
	b := &buffer{}
	b.vbo = gfx.NewVBO()

	layoutData := []struct {
		location gfx.AttributeLocation
		offset   int32
		size     int32
	}{}

	for _, name := range layout {
		var attribute gfx.Attribute
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
					location gfx.AttributeLocation
					offset   int32
					size     int32
				}{location, b.stride, n})
				location++
				b.stride += n
			}
		}
	}

	for _, a := range layoutData {
		vao.SetAttribute(a.location, a.offset, a.size, b.stride, 0, b.vbo)
	}
	return b

}

func (b *buffer) clear() {
	b.data = b.data[:0]
	b.changed = true
}

func (b *buffer) draw(data ...float32) {
	if len(data) != int(b.stride) {
		panic("incorrect vertex size, got " + strconv.Itoa(len(data)) + " expected " + strconv.Itoa(int(b.stride)))
	}
	b.data = append(b.data, data...)
	b.changed = true
}

func (b *buffer) sync() {

	if b.changed {
		if len(b.data) <= b.vboSize {
			b.vbo.SetDataSubsection(0, b.data...)
		} else {
			b.vbo.SetData(b.data...)
			b.vboSize = len(b.data)
		}
		b.changed = false
	}

}

func (b *buffer) count() int32 {
	return int32(len(b.data)) / b.stride
}

func (b *buffer) delete() {
	b.vbo.Delete()
}

type Renderer struct {
	mode     gfx.Mode
	program  gfx.Program
	vao      gfx.VAO
	vertices *buffer
}

func NewRenderer(program gfx.Program, mode gfx.Mode, vertexLayout ...string) *Renderer {
	r := &Renderer{}
	r.mode = mode
	r.program = program
	r.vao = gfx.NewVAO()
	r.vertices = newBuffer(r.vao, r.program.Attributes(), vertexLayout...)
	return r
}

func (r *Renderer) Clear() {
	r.vertices.clear()
}

func (r *Renderer) Draw(vertex ...float32) {
	r.vertices.draw(vertex...)
}

func (r *Renderer) Render() {
	if c := r.vertices.count(); c > 0 {
		r.vertices.sync()
		r.program.Draw(r.vao, r.mode, 0, c)
	}
}

func (r *Renderer) Delete() {
	r.program.Delete()
	r.vao.Delete()
	r.vertices.delete()
}

type InstanceRenderer struct {
	mode      gfx.Mode
	program   gfx.Program
	vao       gfx.VAO
	vertices  *buffer
	instances *buffer
}

func NewInstanceRenderer(program gfx.Program, mode gfx.Mode, vertexLayout []string, instanceLayout ...string) *InstanceRenderer {
	r := &InstanceRenderer{}
	r.mode = mode
	r.program = program
	r.vao = gfx.NewVAO()
	r.vertices = newBuffer(r.vao, r.program.Attributes(), vertexLayout...)
	r.instances = newBuffer(r.vao, r.program.Attributes(), vertexLayout...)
	return r
}

func (r *InstanceRenderer) Clear() {
	r.vertices.clear()
	r.instances.clear()
}

func (r *InstanceRenderer) Draw(vertex ...float32) {
	r.vertices.draw(vertex...)
}

func (r *InstanceRenderer) DrawInstance(instance ...float32) {
	r.instances.draw(instance...)
}

func (r *InstanceRenderer) Render() {
	r.vertices.sync()
	r.instances.sync()
	if vc, ic := r.vertices.count(), r.instances.count(); vc > 0 && ic > 0 {
		r.program.DrawInstanced(r.vao, r.mode, 0, vc, ic)
	}
}

func (r *InstanceRenderer) Delete() {
	r.vao.Delete()
	r.vertices.delete()
	r.instances.delete()
}
