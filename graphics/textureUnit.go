package graphics

import "github.com/go-gl/gl/v4.1-core/gl"

type TextureUnit uint32

func (u TextureUnit) GL() uint32 { return uint32(u + gl.TEXTURE0) }

func (u TextureUnit) SetTexture(t Texture) {
	gl.ActiveTexture(u.GL())
	gl.BindTexture(gl.TEXTURE_2D, t.GL())
}

func (u TextureUnit) WithSetTexture(t Texture) TextureUnit {
	u.SetTexture(t)
	return u
}

func (u TextureUnit) SetTextureArray(t TextureArray) {
	gl.ActiveTexture(u.GL())
	gl.BindTexture(gl.TEXTURE_2D_ARRAY, t.GL())
}

func (u TextureUnit) WithSetTextureArray(t TextureArray) TextureUnit {
	u.SetTextureArray(t)
	return u
}
