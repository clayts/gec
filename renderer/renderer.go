package renderer

import (
	"strconv"

	gfx "github.com/clayts/gec/graphics"
)

type Renderer struct {
	mode    gfx.Mode
	program gfx.Program
	vao     gfx.VAO
	vbo     gfx.VBO
	vboSize int
	stride  int32
	data    []float32
	changed bool
}

func NewRenderer(program gfx.Program, mode gfx.Mode, vertexLayout ...string) *Renderer {
	r := &Renderer{}
	r.mode = mode
	r.program = program
	r.vao = gfx.NewVAO()
	r.vbo = gfx.NewVBO()

	layoutData := []struct {
		location gfx.AttributeLocation
		offset   int32
		size     int32
	}{}

	attributes := r.program.Attributes()
	for _, name := range vertexLayout {
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
				}{location, r.stride, n})
				location++
				r.stride += n
			}
		}
	}

	for _, a := range layoutData {
		r.vao.SetAttribute(a.location, a.offset, a.size, r.stride, r.vbo)
	}
	return r
}

func (r *Renderer) Clear() {
	r.data = r.data[:0]
	r.changed = true
}

func (r *Renderer) Draw(vertex ...float32) {
	if len(vertex) != int(r.stride) {
		panic("incorrect vertex size, got " + strconv.Itoa(len(vertex)) + " expected " + strconv.Itoa(int(r.stride)))
	}
	r.data = append(r.data, vertex...)
	r.changed = true
}

func (r *Renderer) Render() {
	if len(r.data) == 0 {
		return
	}
	if r.changed {
		if len(r.data) <= r.vboSize {
			r.vbo.SetDataSubsection(0, r.data...)
		} else {
			r.vbo.SetData(r.data...)
			r.vboSize = len(r.data)
		}
		r.changed = false
	}
	r.program.Draw(r.vao, r.mode, 0, int32(len(r.data))/r.stride)
}

func (r *Renderer) Delete() {
	r.program.Delete()
	r.vao.Delete()
	r.vbo.Delete()
}
