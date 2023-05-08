package art

import (
	"image"

	geo "github.com/clayts/gec/geometry"
	"github.com/clayts/gec/graphics/automatic"
)

type Sprite struct {
	atlas *Atlas
	index int
}

func (atl *Atlas) MakeSprite(source image.Image) Sprite {
	atl.entries = append(atl.entries, entry{Image: source})
	return Sprite{atlas: atl, index: len(atl.entries) - 1}
}

func (spr Sprite) Image() image.Image { return spr.atlas.entries[spr.index].Image }

func (spr Sprite) Bounds() geo.Rectangle {
	ent := spr.atlas.entries[spr.index]
	return geo.R(geo.V(0, 0), geo.V(float64(ent.w), float64(ent.h)))
}

func (spr Sprite) DrawRegion(
	region geo.Rectangle,
	transform geo.Transform,
	depth float32,
	bufferIndex int,
) {
	ent := spr.atlas.entries[spr.index]
	size := region.Size()
	automatic.Instance{
		float32(transform[0][0]),
		float32(transform[0][1]),
		float32(transform[0][2]),
		float32(transform[1][0]),
		float32(transform[1][1]),
		float32(transform[1][2]),

		depth,

		ent.x + float32(region.Min.X),
		ent.y + float32(region.Min.Y),
		ent.page,
		ent.volume,

		float32(size.X),
		float32(size.Y),
	}.Draw(spr.atlas.buffers[bufferIndex])
}

func (spr Sprite) Draw(
	transform geo.Transform,
	depth float32,
	bufferIndex int,
) {
	ent := spr.atlas.entries[spr.index]
	automatic.Instance{
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
	}.Draw(spr.atlas.buffers[bufferIndex])
}
