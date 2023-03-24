package graphics

import (
	"image"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Texture uint32

func NewTexture(rgba *image.RGBA, smoothing bool, repeat bool, mirror, mipmap bool) Texture {
	var v uint32
	gl.GenTextures(1, &v)
	gl.BindTexture(gl.TEXTURE_2D, v)

	if smoothing {
		if mipmap {
			gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
		} else {
			gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		}
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	} else {
		if mipmap {
			gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST_MIPMAP_NEAREST)
		} else {
			gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
		}
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	}

	if repeat {
		if mirror {
			gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.MIRRORED_REPEAT)
			gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.MIRRORED_REPEAT)
		} else {
			gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
			gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
		}
	} else {
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	}

	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA8,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix),
	)
	if mipmap {
		gl.GenerateMipmap(gl.TEXTURE_2D)
	}

	return Texture(v)
}

func (t Texture) GL() uint32 { return uint32(t) }

func (t Texture) Delete() {
	tgl := t.GL()
	gl.DeleteTextures(1, &tgl)
}
