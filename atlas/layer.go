package atlas

import (
	geo "github.com/clayts/gec/geometry"
	gfx "github.com/clayts/gec/graphics"
)

// TODO get rid of systems

type Layer int

func openLayer() gfx.Buffer {
	vertexLayout := gfx.Layout{
		{program.AttributeLocation("position"), 1, 1, 2},
	}
	instanceLayout := gfx.Layout{
		{program.AttributeLocation("dstTransform"), 1, 2, 3},
		{program.AttributeLocation("dstDepth"), 1, 1, 1},
		{program.AttributeLocation("srcLocation"), 1, 1, 4},
		{program.AttributeLocation("srcSize"), 1, 1, 2},
	}
	vertices := []float32{
		0, 0,
		0, 1,
		1, 0,
		1, 1,
	}
	layer := gfx.OpenBuffer(gfx.TRIANGLE_STRIP, vertexLayout, instanceLayout)
	layer.Vertices().Add(vertices...)
	return layer
}

func (l Layer) Clear() {
	layers[l].Instances().Clear()
}

func (l Layer) Draw(camera geo.Transform) {
	w, h := gfx.Window.GetSize()
	program.SetUniform(program.UniformLocation("screenSize"), [2]float32{float32(w), float32(h)})

	inverse := camera.Inverse()
	program.SetUniform(program.UniformLocation("cameraTransform"), [2][3]float32{
		{float32(inverse[0][0]), float32(inverse[0][1]), float32(inverse[0][2])},
		{float32(inverse[1][0]), float32(inverse[1][1]), float32(inverse[1][2])},
	})

	textureUnits := make([]gfx.TextureUnit, len(volumes))
	for i, volume := range volumes {
		u := gfx.TextureUnit(i)
		u.SetTextureArray(volume)
		textureUnits[i] = u
	}
	program.SetUniform(program.UniformLocation("textureArray"), textureUnits)

	// gl.DepthMask(true)
	// gfx.Clear(false, true, false)

	// gl.Enable(gl.DEPTH_TEST)
	// gl.Disable(gl.BLEND)
	// buf.opaque.Draw(program)

	// gl.DepthMask(false)
	// gl.Enable(gl.BLEND)
	// gl.BlendFunc(gl.ONE, gl.ONE_MINUS_SRC_COLOR)
	// buf.transparent.Draw(program)

	layers[l].Draw(program)
}
