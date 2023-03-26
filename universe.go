package gec

import (
	"time"

	"github.com/clayts/gec/geometry"
	"github.com/clayts/gec/graphics"
	"github.com/clayts/gec/set"
	"github.com/clayts/gec/space"
	"github.com/clayts/gec/sprites"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Universe struct {
	UpdateProcedureSet                  set.Set[func()]
	RenderProcedureSpace                space.Space[func(callShape geometry.Shape)]
	Camera                              geometry.Transform
	OpaqueRenderer, TransparentRenderer *sprites.Renderer
	StepTime                            time.Time
	StepDuration                        time.Duration
}

func NewUniverse() *Universe {
	u := &Universe{}

	u.OpaqueRenderer = sprites.NewRenderer()
	u.TransparentRenderer = sprites.NewRenderer()

	u.Camera = geometry.T()

	u.StepTime = time.Now()

	return u
}

func (u *Universe) Step() {
	now := time.Now()
	u.StepDuration = now.Sub(u.StepTime)
	u.StepTime = now

	u.UpdateProcedureSet.All(func(e *set.Entity[func()]) bool {
		e.Contents()
		return true
	})

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	gl.DepthMask(true)
	gl.Disable(gl.BLEND)
	graphics.Clear(true, true, false)
	u.OpaqueRenderer.Clear()
	u.TransparentRenderer.Clear()

	s := u.Camera.Rectangle(graphics.Bounds())
	u.RenderProcedureSpace.AllIntersecting(s, func(z *space.Zone[func(callShape geometry.Shape)]) bool {
		z.Contents(s)
		return true
	})

	u.OpaqueRenderer.Render()
	gl.DepthMask(false)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.ONE, gl.ONE)
	u.TransparentRenderer.Render()
	graphics.Render()
}

func (u *Universe) Delete() {
	u.OpaqueRenderer.Delete()
	u.TransparentRenderer.Delete()
}
