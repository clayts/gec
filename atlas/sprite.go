package atlas

import (
	"image"

	geo "github.com/clayts/gec/geometry"
)

type Sprite struct {
	index     int
	region    geo.Rectangle
	Layer     Layer
	Transform geo.Transform
	Depth     float32
}

func MakeSprite(source image.Image) Sprite {
	entries = append(entries, entry{Image: source})
	bounds := source.Bounds()
	return Sprite{
		index:     len(entries) - 1,
		region:    geo.R(geo.V(float64(bounds.Min.X), float64(bounds.Min.Y)), geo.V(float64(bounds.Max.X), float64(bounds.Max.Y))),
		Layer:     0,
		Transform: geo.T(),
		Depth:     0,
	}
}

func (spr Sprite) WithRegion(region geo.Rectangle) Sprite {
	// if !spr.region.Contains(region) {
	// 	panic("region out of bounds")
	// }
	spr.region = region
	return spr
}

func (spr Sprite) WithTransform(transform geo.Transform) Sprite {
	spr.Transform = transform
	return spr
}

func (spr Sprite) WithDepth(depth float32) Sprite {
	spr.Depth = depth
	return spr
}

func (spr Sprite) WithLayer(layer Layer) Sprite {
	spr.Layer = layer
	return spr
}

func (spr Sprite) Region() geo.Rectangle { return spr.region }

func (spr Sprite) Image() image.Image { return entries[spr.index] }

func (spr Sprite) Shape() geo.Shape {
	return spr.Transform.Rectangle(geo.R(geo.V(0, 0), spr.region.Size()))
}

func (spr Sprite) Draw() {
	ent := entries[spr.index]
	layer := layers[spr.Layer]
	size := spr.region.Size()
	layer.Instances().Add(
		float32(spr.Transform[0][0]),
		float32(spr.Transform[0][1]),
		float32(spr.Transform[0][2]),
		float32(spr.Transform[1][0]),
		float32(spr.Transform[1][1]),
		float32(spr.Transform[1][2]),

		spr.Depth,

		ent.x+float32(spr.region.Min.X),
		ent.y+float32(spr.region.Min.Y),
		ent.page,
		ent.volume,

		float32(size.X),
		float32(size.Y),
	)
}
