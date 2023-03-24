package graphics

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Mode uint32

const (
	POINTS                   Mode = gl.POINTS
	LINE_STRIP               Mode = gl.LINE_STRIP
	LINE_LOOP                Mode = gl.LINE_LOOP
	LINES                    Mode = gl.LINES
	LINE_STRIP_ADJACENCY     Mode = gl.LINE_STRIP_ADJACENCY
	LINES_ADJACENCY          Mode = gl.LINES_ADJACENCY
	TRIANGLE_STRIP           Mode = gl.TRIANGLE_STRIP
	TRIANGLE_FAN             Mode = gl.TRIANGLE_FAN
	TRIANGLES                Mode = gl.TRIANGLES
	TRIANGLE_STRIP_ADJACENCY Mode = gl.TRIANGLE_STRIP_ADJACENCY
	TRIANGLES_ADJACENCY      Mode = gl.TRIANGLES_ADJACENCY
	PATCHES                  Mode = gl.PATCHES
)

func (m Mode) GL() uint32 { return uint32(m) }
