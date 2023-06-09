package main

import (
	"github.com/clayts/gec/art"
	"github.com/clayts/gec/geometry"
	"github.com/clayts/gec/graphics"
	"github.com/clayts/gec/pixels"
	"github.com/go-gl/gl/v4.1-core/gl"
)

func main() {
	graphics.Open("example")
	defer graphics.Close()

	art.Open()
	defer art.Close()

	const (
		opaque = iota
		transparent
	)
	atlas := art.OpenAtlas(2)
	defer atlas.Close()

	spr := atlas.MakeSprite(pixels.LoadRGBA("test.png"))
	spr2 := atlas.MakeSprite(pixels.LoadRGBA("test2.png"))

	atlas.Pack()

	for !graphics.Window.ShouldClose() {
		spr.Draw(geometry.Translation(geometry.V(10, 10)), 0.5, opaque)
		spr2.Draw(geometry.T(), 0, transparent)

		gl.DepthMask(true)
		graphics.Clear(false, true, false)
		gl.Enable(gl.DEPTH_TEST)
		gl.Disable(gl.BLEND)
		atlas.Buffers(opaque).Draw(geometry.T())

		gl.DepthMask(false)
		gl.Enable(gl.BLEND)
		gl.BlendFunc(gl.ONE, gl.ONE_MINUS_SRC_COLOR)
		atlas.Buffers(transparent).Draw(geometry.T())

		atlas.Buffers(transparent, opaque).Clear()

		graphics.Render()
		graphics.Clear(true, false, false)

	}
}

// package main

// import (
// 	"math/rand"

// 	"github.com/clayts/gec/geometry"
// 	"github.com/clayts/gec/graphics"
// 	"github.com/clayts/gec/images"
// 	"github.com/clayts/gec/sprites"
// 	"github.com/clayts/gec/systems"
// 	"github.com/go-gl/glfw/v3.3/glfw"
// )

// type universe struct {
// 	systems struct {
// 		render *systems.Render
// 		update *systems.Step
// 	}
// }

// func newUniverse() *universe {
// 	u := &universe{}

// 	u.systems.render = systems.NewRender()
// 	u.systems.update = systems.NewStep()

// 	graphics.Window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
// 		switch key {
// 		case glfw.KeyEscape:
// 			graphics.Window.SetShouldClose(true)
// 		}
// 	})

// 	return u
// }

// func (u *universe) createThing(sprite sprites.Sprite, position geometry.Vector, linearVelocity geometry.Vector) {
// 	transform := geometry.Translation(position)
// 	depth := rand.Float32()
// 	shape := sprite.Bounds()

// 	renderZone := u.systems.render.Components.
// 		New(transform.Rectangle(shape)).
// 		SetContents(func(callShape geometry.Shape) {
// 			sprite.Draw(transform, depth)
// 		}).
// 		SetState(true)

// 	u.systems.update.Components.
// 		New().
// 		SetContents(func() {
// 			transform = geometry.Translation(linearVelocity.TimesScalar(u.systems.update.Duration.Seconds())).Times(transform)
// 			renderZone.SetShape(transform.Rectangle(shape))
// 			if !geometry.Contains(u.systems.render.Camera.Rectangle(graphics.Bounds()), renderZone.Shape()) {
// 				linearVelocity = geometry.V(0, 0).Minus(linearVelocity)
// 			}
// 		}).
// 		SetState(true)
// }

// func (u *universe) delete() {
// 	u.systems.render.Delete()
// }

// func main() {
// 	graphics.Initialize("example")
// 	defer graphics.Delete()
// 	u := newUniverse()
// 	defer u.delete()

// 	sprite := u.systems.render.OpaqueRenderer.MakeSprite(images.LoadRGBA("test.png")).SubSprite(geometry.R(geometry.V(0, 50), geometry.V(100, 100)))
// 	for i := 0; i < 100; i++ {
// 		position := geometry.V(
// 			rand.Float64()*(graphics.Bounds().Size().X-sprite.Bounds().Size().X),
// 			rand.Float64()*(graphics.Bounds().Size().Y-sprite.Bounds().Size().Y),
// 		)
// 		linearVelocity := geometry.V(rand.Float64(), rand.Float64()).MinusScalar(0.5).TimesScalar(100)
// 		u.createThing(sprite, position, linearVelocity)
// 	}

// 	sprite2 := u.systems.render.TransparentRenderer.MakeSprite(images.LoadRGBA("test2.png"))
// 	for i := 0; i < 100; i++ {
// 		position := geometry.V(
// 			rand.Float64()*(graphics.Bounds().Size().X-sprite2.Bounds().Size().X),
// 			rand.Float64()*(graphics.Bounds().Size().Y-sprite2.Bounds().Size().Y),
// 		)
// 		linearVelocity := geometry.V(rand.Float64(), rand.Float64()).MinusScalar(0.5).TimesScalar(100)
// 		u.createThing(sprite2, position, linearVelocity)
// 	}

// 	for !graphics.Window.ShouldClose() {
// 		u.systems.update.Step()
// 		u.systems.render.Render()
// 		graphics.Render()
// 		graphics.Clear(true, false, false)
// 	}

// }
