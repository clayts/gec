package graphics

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Mode uint32

const (
	Points                 Mode = gl.POINTS
	LineStrip              Mode = gl.LINE_STRIP
	LineLoop               Mode = gl.LINE_LOOP
	Lines                  Mode = gl.LINES
	LineStripAdjacency     Mode = gl.LINE_STRIP_ADJACENCY
	LinesAdjacency         Mode = gl.LINES_ADJACENCY
	TriangleStrip          Mode = gl.TRIANGLE_STRIP
	TriangleFan            Mode = gl.TRIANGLE_FAN
	Triangles              Mode = gl.TRIANGLES
	TriangleStripAdjacency Mode = gl.TRIANGLE_STRIP_ADJACENCY
	TrianglesAdjacency     Mode = gl.TRIANGLES_ADJACENCY
	Patches                Mode = gl.PATCHES
)

func (m Mode) GL() uint32 { return uint32(m) }
