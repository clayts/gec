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

func (s Sprite) Bounds() geo.Rectangle {
	src := s.sheet.sources[s.index]
	return geo.R(src.offset, geo.V(float64(src.size[0]), float64(src.size[1])))
}

func (s Sprite) Draw(transform geo.Transform, depth float32) {
	s.sheet.initialize()
	src := s.sheet.sources[s.index]
	dst := transform.Times(geo.Translation(src.offset))
	s.sheet.renderer.DrawInstance(
		float32(dst[0][0]), float32(dst[0][1]), float32(dst[0][2]), float32(dst[1][0]), float32(dst[1][1]), float32(dst[1][2]),
		depth,
		src.location[0], src.location[1], src.location[2],
		src.size[0], src.size[1],
	)
}
