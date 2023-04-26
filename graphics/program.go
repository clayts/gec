package graphics

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Program uint32

func OpenProgram(shaders ...Shader) Program {
	program := gl.CreateProgram()
	for _, shader := range shaders {
		gl.AttachShader(program, shader.GL())
	}
	gl.LinkProgram(program)
	return Program(program)
}

func (p Program) GL() uint32 { return uint32(p) }

func (p Program) Close() {
	gl.DeleteProgram(p.GL())
}

func (p Program) Draw(vao VAO, mode Mode, firstVertex, vertexCount int32) {
	gl.UseProgram(p.GL())
	gl.BindVertexArray(vao.GL())
	gl.DrawArrays(mode.GL(), firstVertex, vertexCount)
}

func (p Program) DrawInstanced(vao VAO, mode Mode, firstVertex, vertexCount, instanceCount int32) {
	gl.UseProgram(p.GL())
	gl.BindVertexArray(vao.GL())
	gl.DrawArraysInstanced(mode.GL(), firstVertex, vertexCount, instanceCount)
}

func (p Program) Attributes() []Attribute {
	var attributeCount int32
	gl.GetProgramiv(p.GL(), gl.ACTIVE_ATTRIBUTES, &attributeCount)
	attributes := make([]Attribute, attributeCount)
	for i := int32(0); i < attributeCount; i++ {
		var length int32
		var l int32
		var xtype uint32
		var buf [256]byte
		gl.GetActiveAttrib(p.GL(), uint32(i), int32(len(buf)), &length, &l, &xtype, &buf[0])
		s := gl.GoStr(&buf[0])
		locationUnchecked := gl.GetAttribLocation(p.GL(), gl.Str(s+"\x00"))
		if locationUnchecked < 0 {
			panic("invalid attribute " + s)
		}
		location := AttributeLocation(locationUnchecked)
		m, n := dimensions(xtype)
		attributes[i].name = s
		attributes[i].location = location
		attributes[i].l = l
		attributes[i].m = m
		attributes[i].n = n
	}
	return attributes
}

func dimensions(t uint32) (m, n int32) {
	switch t {
	case gl.FLOAT:
		return 1, 1
	case gl.FLOAT_VEC2:
		return 1, 2
	case gl.FLOAT_VEC3:
		return 1, 3
	case gl.FLOAT_VEC4:
		return 1, 4
	case gl.FLOAT_MAT2:
		return 2, 2
	case gl.FLOAT_MAT3:
		return 3, 3
	case gl.FLOAT_MAT4:
		return 4, 4
	case gl.FLOAT_MAT2x3:
		return 2, 3
	case gl.FLOAT_MAT2x4:
		return 2, 4
	case gl.FLOAT_MAT3x2:
		return 3, 2
	case gl.FLOAT_MAT3x4:
		return 3, 4
	case gl.FLOAT_MAT4x2:
		return 4, 2
	case gl.FLOAT_MAT4x3:
		return 4, 3
	}
	panic("invalid attribute type")
}

func (p Program) UniformLocation(uniformName string) UniformLocation {
	return UniformLocation(gl.GetUniformLocation(p.GL(), gl.Str(uniformName+"\x00")))
}

func (p Program) SetUniform(u UniformLocation, data interface{}) {
	gl.UseProgram(p.GL())
	switch data := data.(type) {

	case []TextureUnit:
		slice := make([]int32, len(data))
		for i, tu := range data {
			slice[i] = int32(tu)
		}
		gl.Uniform1iv(u.GL(), int32(len(slice)), &slice[0])

	case TextureUnit:
		gl.Uniform1i(u.GL(), int32(data)) //NOT data.GL() for some reason (requires e.g. 0 for gl.TEXTURE0)

	case float32:
		gl.Uniform1f(u.GL(), data)

	case [2]float32:
		gl.Uniform2f(u.GL(), data[0], data[1])
	case [3]float32:
		gl.Uniform3f(u.GL(), data[0], data[1], data[2])
	case [4]float32:
		gl.Uniform4f(u.GL(), data[0], data[1], data[2], data[3])

	case [2][2]float32:
		buf := [4]float32{data[0][0], data[0][1], data[1][0], data[1][1]}
		gl.UniformMatrix2fv(u.GL(), 1, false, &buf[0])
	case [2][3]float32:
		buf := [6]float32{data[0][0], data[0][1], data[0][2], data[1][0], data[1][1], data[1][2]}
		gl.UniformMatrix2x3fv(u.GL(), 1, false, &buf[0])
	case [2][4]float32:
		buf := [8]float32{data[0][0], data[0][1], data[0][2], data[0][3], data[1][0], data[1][1], data[1][2], data[1][3]}
		gl.UniformMatrix2x4fv(u.GL(), 1, false, &buf[0])

	case [3][2]float32:
		buf := [6]float32{data[0][0], data[0][1], data[1][0], data[1][1], data[2][0], data[2][1]}
		gl.UniformMatrix3x2fv(u.GL(), 1, false, &buf[0])
	case [3][3]float32:
		buf := [9]float32{data[0][0], data[0][1], data[0][2], data[1][0], data[1][1], data[1][2], data[2][0], data[2][1], data[2][2]}
		gl.UniformMatrix3fv(u.GL(), 1, false, &buf[0])
	case [3][4]float32:
		buf := [12]float32{data[0][0], data[0][1], data[0][2], data[0][3], data[1][0], data[1][1], data[1][2], data[1][3], data[2][0], data[2][1], data[2][2], data[2][3]}
		gl.UniformMatrix3x4fv(u.GL(), 1, false, &buf[0])

	case [4][2]float32:
		buf := [8]float32{data[0][0], data[0][1], data[1][0], data[1][1], data[2][0], data[2][1], data[3][0], data[3][1]}
		gl.UniformMatrix4x2fv(u.GL(), 1, false, &buf[0])
	case [4][3]float32:
		buf := [12]float32{data[0][0], data[0][1], data[0][2], data[1][0], data[1][1], data[1][2], data[2][0], data[2][1], data[2][2], data[3][0], data[3][1], data[3][2]}
		gl.UniformMatrix4x3fv(u.GL(), 1, false, &buf[0])
	case [4][4]float32:
		buf := [16]float32{data[0][0], data[0][1], data[0][2], data[0][3], data[1][0], data[1][1], data[1][2], data[1][3], data[2][0], data[2][1], data[2][2], data[2][3], data[3][0], data[3][1], data[3][2], data[3][3]}
		gl.UniformMatrix4fv(u.GL(), 1, false, &buf[0])

	default:
		panic("unsupported data type")
	}
}
