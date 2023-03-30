package sprites

import (
	gfx "github.com/clayts/gec/graphics"

	_ "embed"
)

var (
	//go:embed shaders/vertex.glsl
	vertexShaderSource string

	//go:embed shaders/fragment.glsl
	fragmentShaderSource string

	vertexShader, fragmentShader gfx.Shader
	program                      gfx.Program

	cameraTransformUniformLocation, screenSizeUniformLocation, textureArrayUniformLocation gfx.UniformLocation
)

func Initialize() {
	vertexShader = gfx.NewVertexShader(vertexShaderSource)
	fragmentShader = gfx.NewFragmentShader(fragmentShaderSource)
	program = gfx.NewProgram(vertexShader, fragmentShader)
	cameraTransformUniformLocation = program.UniformLocation("cameraTransform")
	screenSizeUniformLocation = program.UniformLocation("screenSize")
	textureArrayUniformLocation = program.UniformLocation("textureArray")

}

func Delete() {
	program.Delete()
	vertexShader.Delete()
	fragmentShader.Delete()
}
