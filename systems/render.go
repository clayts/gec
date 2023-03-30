package systems

import (
	"github.com/clayts/gec/geometry"
	"github.com/clayts/gec/graphics"
	"github.com/clayts/gec/space"
	"github.com/clayts/gec/sprites"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Render struct {
	Components space.Space[func(callShape geometry.Shape)]
	Camera     geometry.Transform
	OpaqueRenderer,
	TransparentRenderer *sprites.Sheet
}

func NewRender() *Render {
	r := &Render{}

	r.OpaqueRenderer = sprites.NewSheet()
	r.TransparentRenderer = sprites.NewSheet()

	r.Camera = geometry.T()

	return r
}

func (r *Render) Render() {
	r.OpaqueRenderer.Clear()
	r.TransparentRenderer.Clear()

	s := r.Camera.Rectangle(graphics.Bounds())
	r.Components.AllIntersecting(s, func(z *space.Zone[func(callShape geometry.Shape)]) bool {
		z.Contents(s)
		return true
	})

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	gl.DepthMask(true)
	gl.Disable(gl.BLEND)
	graphics.Clear(true, true, false)
	r.OpaqueRenderer.Render(r.Camera)
	gl.DepthMask(false)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.ONE, gl.ONE_MINUS_SRC_COLOR)
	r.TransparentRenderer.Render(r.Camera)
	graphics.Render()
}

func (r *Render) Delete() {
	r.OpaqueRenderer.Delete()
	r.TransparentRenderer.Delete()
}
