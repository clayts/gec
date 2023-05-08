package graphics

import (
	"github.com/clayts/gec/geometry"
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

func (p Program) AttributeLocation(name string) AttributeLocation {
	return AttributeLocation(gl.GetAttribLocation(p.GL(), gl.Str(name+"\x00")))
}

func (p Program) UniformLocation(uniformName string) UniformLocation {
	return UniformLocation(gl.GetUniformLocation(p.GL(), gl.Str(uniformName+"\x00")))
}

func (p Program) SetUniform(u UniformLocation, data interface{}) {
	gl.UseProgram(p.GL())
	switch data := data.(type) {

	case []TextureUnit:
		if len(data) > 0 {
			slice := make([]int32, len(data))
			for i, tu := range data {
				slice[i] = int32(tu)
			}
			gl.Uniform1iv(u.GL(), int32(len(slice)), &slice[0])
		}

	case TextureUnit:
		gl.Uniform1i(u.GL(), int32(data)) //NOT data.GL() for some reason (requires e.g. 0 for gl.TEXTURE0)

	case geometry.Transform:
		buf := [6]float32{float32(data[0][0]), float32(data[0][1]), float32(data[0][2]), float32(data[1][0]), float32(data[1][1]), float32(data[1][2])}
		gl.UniformMatrix2x3fv(u.GL(), 1, false, &buf[0])

	case geometry.Vector:
		gl.Uniform2f(u.GL(), float32(data.X), float32(data.Y))

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
