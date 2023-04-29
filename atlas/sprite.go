package atlas

import (
	"image"

	geo "github.com/clayts/gec/geometry"
)

type Sprite struct {
	index        int
	region       geo.Rectangle
	Transform    geo.Transform
	Depth        float32
	Transparency bool
}

func MakeSprite(source image.Image) Sprite {
	entries = append(entries, entry{Image: source})
	bounds := source.Bounds()
	return Sprite{
		index:        len(entries) - 1,
		region:       geo.R(geo.V(float64(bounds.Min.X), float64(bounds.Min.Y)), geo.V(float64(bounds.Max.X), float64(bounds.Max.Y))),
		Transform:    geo.T(),
		Depth:        0,
		Transparency: false,
	}
}

func (spr Sprite) WithRegion(region geo.Rectangle) Sprite {
	if !spr.region.Contains(region) {
		panic("region out of bounds")
	}
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

func (spr Sprite) WithTransparency(transparency bool) Sprite {
	spr.Transparency = transparency
	return spr
}

func (spr Sprite) Region() geo.Rectangle { return spr.region }

func (spr Sprite) Image() image.Image { return entries[spr.index] }

func (spr Sprite) Shape() geo.Shape {
	return spr.Transform.Rectangle(geo.R(geo.V(0, 0), spr.region.Size()))
}
