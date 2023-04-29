package graphics

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

type VBO uint32

func OpenVBO(data ...float32) VBO {
	var vgl uint32
	gl.GenBuffers(1, &vgl)

	v := VBO(vgl)
	if len(data) > 0 {
		v.SetData(data...)
	}
	return v
}

func (v VBO) GL() uint32 { return uint32(v) }

func (v VBO) Close() {
	vgl := v.GL()
	gl.DeleteBuffers(1, &vgl)
}

func (v VBO) SetData(data ...float32) {
	gl.BindBuffer(gl.ARRAY_BUFFER, v.GL())
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(data), gl.Ptr(data), gl.STATIC_DRAW)
}

func (v VBO) SetDataSubsection(offset int, data ...float32) {
	gl.BindBuffer(gl.ARRAY_BUFFER, v.GL())
	gl.BufferSubData(gl.ARRAY_BUFFER, 4*offset, 4*len(data), gl.Ptr(data))
}

func (v VBO) Type() uint32 { return gl.FLOAT }

func (v VBO) Normalized() bool { return false }
