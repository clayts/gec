package graphics

import (
	"image"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type TextureArray uint32

// Images must be the same size
func NewTextureArray(images []*image.RGBA, smoothing, repeat, mirror, mipmap bool) TextureArray {
	width := int32(images[0].Bounds().Dx())
	height := int32(images[0].Bounds().Dy())
	var v uint32
	gl.GenTextures(1, &v)
	gl.BindTexture(gl.TEXTURE_2D_ARRAY, v)

	if smoothing {
		if mipmap {
			gl.TexParameteri(gl.TEXTURE_2D_ARRAY, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
		} else {
			gl.TexParameteri(gl.TEXTURE_2D_ARRAY, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		}
		gl.TexParameteri(gl.TEXTURE_2D_ARRAY, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	} else {
		if mipmap {
			gl.TexParameteri(gl.TEXTURE_2D_ARRAY, gl.TEXTURE_MIN_FILTER, gl.NEAREST_MIPMAP_NEAREST)
		} else {
			gl.TexParameteri(gl.TEXTURE_2D_ARRAY, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
		}
		gl.TexParameteri(gl.TEXTURE_2D_ARRAY, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	}

	if repeat {
		if mirror {
			gl.TexParameteri(gl.TEXTURE_2D_ARRAY, gl.TEXTURE_WRAP_S, gl.MIRRORED_REPEAT)
			gl.TexParameteri(gl.TEXTURE_2D_ARRAY, gl.TEXTURE_WRAP_T, gl.MIRRORED_REPEAT)
		} else {
			gl.TexParameteri(gl.TEXTURE_2D_ARRAY, gl.TEXTURE_WRAP_S, gl.REPEAT)
			gl.TexParameteri(gl.TEXTURE_2D_ARRAY, gl.TEXTURE_WRAP_T, gl.REPEAT)
		}
	} else {
		gl.TexParameteri(gl.TEXTURE_2D_ARRAY, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_2D_ARRAY, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	}

	gl.TexStorage3D(gl.TEXTURE_2D_ARRAY, 1, gl.RGBA8, width, height, int32(len(images)))
	for i, img := range images {
		gl.TexSubImage3D(gl.TEXTURE_2D_ARRAY, 0, 0, 0, int32(i), width, height, 1, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))
	}
	if mipmap {
		gl.GenerateMipmap(gl.TEXTURE_2D_ARRAY)
	}

	return TextureArray(v)
}

func (t TextureArray) GL() uint32 { return uint32(t) }

func (t TextureArray) Delete() {
	v := t.GL()
	gl.DeleteTextures(1, &v)
}
