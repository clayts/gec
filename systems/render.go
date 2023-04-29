package systems

import (
	"github.com/clayts/gec/atlas"
	"github.com/clayts/gec/geometry"
	"github.com/clayts/gec/graphics"
	"github.com/clayts/gec/space"
)

type Render struct {
	Components space.Space[func(callShape geometry.Shape)]
	Camera     geometry.Transform
	Buffer     atlas.Atlas
}

func OpenRender() *Render {
	r := &Render{}

	r.Buffer = atlas.OpenAtlas()

	r.Camera = geometry.T()

	return r
}

func (r *Render) Render() {

	s := r.Camera.Rectangle(graphics.Bounds())
	r.Components.AllIntersecting(s, func(z *space.Zone[func(callShape geometry.Shape)]) bool {
		z.Contents(s)
		return true
	})

	r.Buffer.Draw(r.Camera)
}

func (r *Render) Close() {
	r.Buffer.Close()
}
