package atlas

import (
	geo "github.com/clayts/gec/geometry"
	gfx "github.com/clayts/gec/graphics"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Buffer struct {
	opaque      gfx.Buffer
	transparent gfx.Buffer
}

func OpenBuffer() Buffer {
	buf := Buffer{}
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
	buf.opaque = gfx.OpenBuffer(gfx.TRIANGLE_STRIP, vertexLayout, instanceLayout)
	buf.opaque.Vertices().Add(vertices...)
	buf.transparent = gfx.OpenBuffer(gfx.TRIANGLE_STRIP, vertexLayout, instanceLayout)
	buf.transparent.Vertices().Add(vertices...)
	return buf
}

func (buf Buffer) Draw(camera geo.Transform) {
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

	gl.DepthMask(true)
	gfx.Clear(false, true, false)

	gl.Enable(gl.DEPTH_TEST)
	gl.Disable(gl.BLEND)
	buf.opaque.Draw(program)

	gl.DepthMask(false)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.ONE, gl.ONE_MINUS_SRC_COLOR)
	buf.transparent.Draw(program)

}

func (buf Buffer) Add(spr Sprite) {
	ent := entries[spr.index]
	sub := buf.opaque
	if spr.Transparency {
		sub = buf.transparent
	}
	size := spr.region.Size()
	sub.Instances().Add(
		float32(spr.Transform[0][0]), float32(spr.Transform[0][1]), float32(spr.Transform[0][2]), float32(spr.Transform[1][0]), float32(spr.Transform[1][1]), float32(spr.Transform[1][2]),
		spr.Depth,
		ent.x+float32(spr.region.Min.X), ent.y+float32(spr.region.Min.Y), ent.page, ent.volume,
		float32(size.X), float32(size.Y),
	)
}

func (buf Buffer) Clear() {
	buf.transparent.Instances().Clear()
	buf.opaque.Instances().Clear()
}

func (buf Buffer) Close() {
	buf.opaque.Close()
	buf.transparent.Close()
}
