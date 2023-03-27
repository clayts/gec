package main

import (
	"math/rand"

	"github.com/clayts/gec"
	"github.com/clayts/gec/geometry"
	"github.com/clayts/gec/graphics"
	"github.com/clayts/gec/images"
	"github.com/clayts/gec/sprites"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type universe struct {
	*gec.Universe
}

func newUniverse() *universe {
	u := &universe{}

	u.Universe = gec.NewUniverse()

	graphics.Window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		switch key {
		case glfw.KeyEscape:
			graphics.Window.SetShouldClose(true)
		}
	})

	return u
}

func (u *universe) createThing(sprite sprites.Sprite, position geometry.Vector, linearVelocity geometry.Vector) {
	transform := geometry.Translation(position)
	depth := rand.Float32()

	renderZone := u.RenderProcedureSpace.NewZone().
		SetContents(func(callShape geometry.Shape) {
			sprite.Draw(transform, depth)
		}).
		SetShape(transform.Rectangle(sprite.Bounds())).
		Enable()

	u.UpdateProcedureSet.NewEntity().
		SetContents(func() {
			transform = geometry.Translation(linearVelocity.TimesScalar(u.StepDuration.Seconds())).Times(transform)
			renderZone.SetShape(transform.Rectangle(sprite.Bounds()))
			if !geometry.Contains(u.Camera.Rectangle(graphics.Bounds()), renderZone.Shape()) {
				linearVelocity = geometry.V(0, 0).Minus(linearVelocity)
			}
		}).
		Enable()
}

func main() {
	graphics.Initialize("example")
	defer graphics.Delete()

	sprites.Initialize()
	defer sprites.Delete()

	u := newUniverse()
	defer u.Delete()

	sprite := u.OpaqueRenderer.MakeSprite(images.LoadRGBA("test.png"))
	for i := 0; i < 100; i++ {
		position := geometry.V(
			rand.Float64()*(graphics.Bounds().Size().X-sprite.Bounds().Size().X),
			rand.Float64()*(graphics.Bounds().Size().Y-sprite.Bounds().Size().Y),
		)
		linearVelocity := geometry.V(rand.Float64()*100, rand.Float64()*100)
		u.createThing(sprite, position, linearVelocity)
	}

	sprite2 := u.TransparentRenderer.MakeSprite(images.LoadRGBA("test2.png"))
	for i := 0; i < 100; i++ {
		position := geometry.V(
			rand.Float64()*(graphics.Bounds().Size().X-sprite2.Bounds().Size().X),
			rand.Float64()*(graphics.Bounds().Size().Y-sprite2.Bounds().Size().Y),
		)
		linearVelocity := geometry.V(rand.Float64(), rand.Float64()).MinusScalar(0.5).TimesScalar(100)
		u.createThing(sprite2, position, linearVelocity)
	}

	for !graphics.Window.ShouldClose() {
		u.Step()
	}

}
