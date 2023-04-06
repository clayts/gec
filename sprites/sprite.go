package sprites

import (
	"image"

	geo "github.com/clayts/gec/geometry"
)

type Sprite struct {
	sheet *Sheet
	index int
}

func (sh *Sheet) MakeSprite(img image.Image) Sprite {
	sh.sources = append(sh.sources, struct {
		location [3]float32
		size     [2]float32
		offset   geo.Vector
		image    image.Image
	}{
		image:  img,
		size:   [2]float32{float32(img.Bounds().Dx()), float32(img.Bounds().Dy())},
		offset: geo.V(float64(img.Bounds().Min.X), float64(img.Bounds().Min.Y)),
	})
	return Sprite{
		sheet: sh,
		index: len(sh.sources) - 1,
	}
}

func (sh *Sheet) MakeSprites(imgs ...image.Image) []Sprite {
	ss := make([]Sprite, len(imgs))
	for i, img := range imgs {
		ss[i] = sh.MakeSprite(img)
	}
	return ss
}

func (s Sprite) Size() geo.Vector {
	src := s.sheet.sources[s.index]
	return geo.V(float64(src.size[0]), float64(src.size[1]))
}

func (s Sprite) Offset() geo.Vector {
	src := s.sheet.sources[s.index]
	return src.offset
}

func (s Sprite) Instance(transform geo.Transform, depth float32) Instance {
	return Instance{s, transform, depth}
}

type Instance struct {
	Sprite    Sprite
	Transform geo.Transform
	Depth     float32
}

func (i Instance) Draw() {
	i.Sprite.sheet.initialize()
	src := i.Sprite.sheet.sources[i.Sprite.index]
	dst := i.Transform.Times(geo.Translation(src.offset))
	i.Sprite.sheet.renderer.DrawInstance(
		float32(dst[0][0]), float32(dst[0][1]), float32(dst[0][2]), float32(dst[1][0]), float32(dst[1][1]), float32(dst[1][2]),
		i.Depth,
		src.location[0], src.location[1], src.location[2],
		src.size[0], src.size[1],
	)
}

func (i Instance) Shape() geo.Shape {
	return i.Transform.Rectangle(geo.R(i.Sprite.Offset(), i.Sprite.Offset().Plus(i.Sprite.Size())))
}
