package graphics

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

type VAO uint32

func NewVAO() VAO {
	var v uint32

	gl.GenVertexArrays(1, &v)

	return VAO(v)
}

func (v VAO) SetAttribute(a AttributeLocation, offset, size, stride int32, divisor uint32, vbo VBO) {
	gl.BindVertexArray(v.GL())
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo.GL())
	gl.EnableVertexAttribArray(a.GL())
	gl.VertexAttribPointerWithOffset(a.GL(), size, vbo.Type(), vbo.Normalized(), stride*4, uintptr(offset*4))
	gl.VertexAttribDivisor(a.GL(), divisor)
}

func (v VAO) GL() uint32 { return uint32(v) }

func (v VAO) Delete() {
	vgl := v.GL()
	gl.DeleteVertexArrays(1, &vgl)
}
