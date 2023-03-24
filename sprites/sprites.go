package sprites

import (
	gfx "github.com/clayts/gec/graphics"

	_ "embed"
)

var (
	//go:embed shaders/vertex.glsl
	vertexShaderSource string

	//go:embed shaders/geometry.glsl
	geometryShaderSource string

	//go:embed shaders/fragment.glsl
	fragmentShaderSource string

	vertexShader, geometryShader, fragmentShader gfx.Shader
	program                                      gfx.Program

	screenSizeUniformLocation, textureArrayUniformLocation gfx.UniformLocation
)

func Initialize() {
	vertexShader = gfx.NewVertexShader(vertexShaderSource)
	geometryShader = gfx.NewGeometryShader(geometryShaderSource)
	fragmentShader = gfx.NewFragmentShader(fragmentShaderSource)
	program = gfx.NewProgram(vertexShader, geometryShader, fragmentShader)
	screenSizeUniformLocation = program.UniformLocation("screenSize")
	textureArrayUniformLocation = program.UniformLocation("textureArray")
}

func Delete() {
	program.Delete()
	vertexShader.Delete()
	geometryShader.Delete()
	fragmentShader.Delete()
}
