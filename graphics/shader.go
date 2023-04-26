package graphics

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Shader uint32

func OpenVertexShader(source string) Shader {
	return newShader(source, gl.VERTEX_SHADER)
}

func OpenGeometryShader(source string) Shader {
	return newShader(source, gl.GEOMETRY_SHADER)
}

func OpenFragmentShader(source string) Shader {
	return newShader(source, gl.FRAGMENT_SHADER)
}

func (sh Shader) GL() uint32 { return uint32(sh) }

func (sh Shader) Delete() {
	gl.DeleteShader(uint32(sh))
}

func newShader(source string, shaderType uint32) Shader {
	v, err := compileShader(source, shaderType)
	if err != nil {
		panic(err)
	}
	return Shader(v)
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source + "\x00")
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}
