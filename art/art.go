package art

import (
	_ "embed"
	"image"

	gfx "github.com/clayts/gec/graphics"
)

var (
	//go:embed shaders/vertex.glsl
	vertexShaderSource string

	//go:embed shaders/fragment.glsl
	fragmentShaderSource string

	shaders []gfx.Shader
	program gfx.Program
)

type entry struct {
	image.Image
	volume, page, x, y, w, h float32
}

func Open() {
	shaders = []gfx.Shader{
		gfx.OpenVertexShader(vertexShaderSource),
		gfx.OpenFragmentShader(fragmentShaderSource),
	}
	program = gfx.OpenProgram(shaders...)
}

func Close() {
	program.Close()

	for _, shader := range shaders {
		shader.Close()
	}
	shaders = shaders[:0]

}
