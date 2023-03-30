package graphics

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Program uint32

func NewProgram(shaders ...Shader) Program {
	program := gl.CreateProgram()
	for _, shader := range shaders {
		gl.AttachShader(program, shader.GL())
	}
	gl.LinkProgram(program)
	return Program(program)
}

func (p Program) GL() uint32 { return uint32(p) }

func (p Program) Delete() {
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

func (p *Program) SetUniform(u UniformLocation, data interface{}) {
	gl.UseProgram(p.GL())
	switch data := data.(type) {
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

	case int32:
		gl.Uniform1i(u.GL(), data)
	case [2]int32:
		gl.Uniform2i(u.GL(), data[0], data[1])
	case [3]int32:
		gl.Uniform3i(u.GL(), data[0], data[1], data[2])
	case [4]int32:
		gl.Uniform4i(u.GL(), data[0], data[1], data[2], data[3])

	case uint32:
		gl.Uniform1ui(u.GL(), data)
	case [2]uint32:
		gl.Uniform2ui(u.GL(), data[0], data[1])
	case [3]uint32:
		gl.Uniform3ui(u.GL(), data[0], data[1], data[2])
	case [4]uint32:
		gl.Uniform4ui(u.GL(), data[0], data[1], data[2], data[3])

	default:
		panic("unsupported data type")
	}
}
