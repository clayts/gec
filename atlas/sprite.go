package atlas

import (
	"image"

	geo "github.com/clayts/gec/geometry"
)

type Sprite int

func MakeSprite(source image.Image) Sprite {
	entries = append(entries, entry{Image: source})
	return Sprite(len(entries) - 1)
}

func (spr Sprite) Image() image.Image { return entries[spr].Image }

func (spr Sprite) Bounds() geo.Rectangle {
	ent := entries[spr]
	return geo.R(geo.V(0, 0), geo.V(float64(ent.w), float64(ent.h)))
}

func (spr Sprite) DrawRegion(
	region geo.Rectangle,
	transform geo.Transform,
	depth float32,
	buffer Buffer,
) {
	ent := entries[spr]
	size := region.Size()
	buffer.internal.Instances().Add(
		float32(transform[0][0]),
		float32(transform[0][1]),
		float32(transform[0][2]),
		float32(transform[1][0]),
		float32(transform[1][1]),
		float32(transform[1][2]),

		depth,

		ent.x+float32(region.Min.X),
		ent.y+float32(region.Min.Y),
		ent.page,
		ent.volume,

		float32(size.X),
		float32(size.Y),
	)
}

func (spr Sprite) Draw(
	transform geo.Transform,
	depth float32,
	buffer Buffer,
) {
	ent := entries[spr]
	buffer.internal.Instances().Add(
		float32(transform[0][0]),
		float32(transform[0][1]),
		float32(transform[0][2]),
		float32(transform[1][0]),
		float32(transform[1][1]),
		float32(transform[1][2]),

		depth,

		ent.x,
		ent.y,
		ent.page,
		ent.volume,

		ent.w,
		ent.h,
	)
}
