package art

import (
	gfx "github.com/clayts/gec/graphics"
	auto "github.com/clayts/gec/graphics/automatic"
)

func openBuffer() auto.Buffer {
	vertexLayout := auto.Layout{
		{program.AttributeLocation("position"), 1, 1, 2},
	}
	instanceLayout := auto.Layout{
		{program.AttributeLocation("dstTransform"), 1, 2, 3},
		{program.AttributeLocation("dstDepth"), 1, 1, 1},
		{program.AttributeLocation("srcLocation"), 1, 1, 4},
		{program.AttributeLocation("srcSize"), 1, 1, 2},
	}
	vertices := auto.Vertex{
		0, 0,
		0, 1,
		1, 0,
		1, 1,
	}
	buf := auto.OpenBuffer(gfx.TriangleStrip, vertexLayout, instanceLayout)
	vertices.Draw(buf)
	return buf
}
