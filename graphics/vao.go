package graphics

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

type VAO uint32

func OpenVAO() VAO {
	var v uint32

	gl.GenVertexArrays(1, &v)

	return VAO(v)
}

func (vao VAO) Draw(p Program, mode Mode, firstVertex, vertexCount int32) {
	gl.UseProgram(p.GL())
	gl.BindVertexArray(vao.GL())
	gl.DrawArrays(mode.GL(), firstVertex, vertexCount)
}

func (vao VAO) DrawInstanced(p Program, mode Mode, firstVertex, vertexCount, instanceCount int32) {
	gl.UseProgram(p.GL())
	gl.BindVertexArray(vao.GL())
	gl.DrawArraysInstanced(mode.GL(), firstVertex, vertexCount, instanceCount)
}

func (v VAO) SetAttribute(a AttributeLocation, offset, size, stride int32, divisor uint32, vbo VBO) {
	gl.BindVertexArray(v.GL())
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo.GL())
	gl.EnableVertexAttribArray(a.GL())
	gl.VertexAttribPointerWithOffset(a.GL(), size, vbo.Type(), vbo.Normalized(), stride*4, uintptr(offset*4))
	gl.VertexAttribDivisor(a.GL(), divisor)
}

func (v VAO) GL() uint32 { return uint32(v) }

func (v VAO) Close() {
	vgl := v.GL()
	gl.DeleteVertexArrays(1, &vgl)
}
