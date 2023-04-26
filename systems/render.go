package systems

import (
	"github.com/clayts/gec/atlas"
	"github.com/clayts/gec/geometry"
	"github.com/clayts/gec/graphics"
	"github.com/clayts/gec/space"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Render struct {
	Components           space.Space[func(callShape geometry.Shape)]
	Camera               geometry.Transform
	Renderer             atlas.Renderer
	OpaqueInstances      []float32
	TransparentInstances []float32
}

func NewRender() *Render {
	r := &Render{}

	r.Renderer = atlas.OpenRenderer()

	r.Camera = geometry.T()

	return r
}

func (r *Render) Render() {

	s := r.Camera.Rectangle(graphics.Bounds())
	r.Components.AllIntersecting(s, func(z *space.Zone[func(callShape geometry.Shape)]) bool {
		z.Contents(s)
		return true
	})

	gl.DepthMask(true)
	graphics.Clear(false, true, false)

	gl.Enable(gl.DEPTH_TEST)
	gl.Disable(gl.BLEND)
	r.Renderer.SetInstances(r.OpaqueInstances...)
	r.Renderer.Render(r.Camera)

	gl.DepthMask(false)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.ONE, gl.ONE_MINUS_SRC_COLOR)
	r.Renderer.SetInstances(r.TransparentInstances...)
	r.Renderer.Render(r.Camera)
}

func (r *Render) Delete() {
	r.Renderer.Close()
}
