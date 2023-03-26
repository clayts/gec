package main

import (
	"fmt"
	"math/rand"

	"github.com/clayts/gec/geometry"
	"github.com/clayts/gec/graphics"
	"github.com/clayts/gec/image"
	"github.com/clayts/gec/space"
	"github.com/clayts/gec/sprites"
	"github.com/go-gl/gl/v4.1-core/gl"
)

func main() {
	spc := space.New[struct{}]()
	leaves := []*space.Zone[struct{}]{}
	for i := 0; i < 100000; i++ {
		v := geometry.V(rand.Float64()*100000, rand.Float64()*100000)
		x := spc.NewLeaf().SetShape(geometry.R(v, v.Plus(geometry.V(10, 10)))).Enable()
		leaves = append(leaves, x)
	}
	fmt.Println("-------------")
	spc.AllZonesIntersecting(geometry.R(geometry.V(0, 0), geometry.V(1000, 1000)), func(l *space.Zone[struct{}]) bool {
		fmt.Println("found", l)
		return true
	})
	fmt.Println("-------------")
	for _, l := range leaves {
		l.Disable()
	}
	fmt.Println(spc)

	// ------------------------------------------------------

	graphics.Initialize("test")
	defer graphics.Delete()

	sprites.Initialize()
	defer sprites.Delete()

	r := sprites.NewRenderer()
	defer r.Delete()

	s := r.NewSprite(image.LoadRGBA("test.png"))
	ts := make([]geometry.Transform, 10000)
	ds := make([]float32, len(ts))
	for i := range ts {
		ts[i] = geometry.Translation(geometry.V(rand.Float64()*100000, rand.Float64()*50000))
		ds[i] = rand.Float32()
	}

	s2 := r.NewSprite(image.LoadRGBA("test2.png"))
	ts2 := make([]geometry.Transform, 10000)
	ds2 := make([]float32, len(ts))
	for i := range ts2 {
		ts2[i] = geometry.Translation(geometry.V(rand.Float64()*100000, rand.Float64()*50000))
		ds2[i] = rand.Float32()
	}

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	for !graphics.Window.ShouldClose() {
		dt := graphics.DeltaTime()
		fmt.Println(1 / dt)
		gl.DepthMask(true)
		gl.Disable(gl.BLEND)
		graphics.Clear(true, true, false)
		r.Clear()
		for i, t := range ts {
			ts[i] = geometry.RotationAround(geometry.Angle(36*dt), t.Vector(s.Bounds().Center())).Times(t)
			s.Draw(t, ds[i])
		}
		r.Render()
		r.Clear()
		for i, t := range ts2 {
			ts2[i] = geometry.RotationAround(geometry.Angle(36*dt), t.Vector(s2.Bounds().Center())).Times(t)
			s2.Draw(t, ds2[i])
		}
		gl.DepthMask(false)
		gl.Enable(gl.BLEND)
		gl.BlendFunc(gl.ONE, gl.ONE)
		r.Render()
		graphics.Render()
	}

}
