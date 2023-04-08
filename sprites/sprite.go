package sprites

import (
	"image"

	geo "github.com/clayts/gec/geometry"
)

type Sprite struct {
	sheet      *Sheet
	imageIndex int
	min        [2]float32
	size       [2]float32
}

func (sh *Sheet) MakeSprite(img image.Image) Sprite {
	sh.images = append(sh.images, struct {
		location [3]float32
		image    image.Image
	}{
		image: img,
	})
	return Sprite{
		sheet:      sh,
		imageIndex: len(sh.images) - 1,
		size:       [2]float32{float32(img.Bounds().Dx()), float32(img.Bounds().Dy())},
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
	return geo.R(geo.V(0, 0), geo.V(float64(s.size[0]), float64(s.size[1])))
}

func (s Sprite) Draw(transform geo.Transform, depth float32) {
	s.sheet.initialize()
	src := s.sheet.images[s.imageIndex]
	s.sheet.renderer.DrawInstance(
		float32(transform[0][0]), float32(transform[0][1]), float32(transform[0][2]), float32(transform[1][0]), float32(transform[1][1]), float32(transform[1][2]),
		depth,
		src.location[0]+s.min[0], src.location[1]+s.min[1], src.location[2],
		s.size[0], s.size[1],
	)
}

func (s Sprite) SubSprite(region geo.Rectangle) Sprite {
	if !s.Bounds().Contains(region) {
		panic("region out of bounds")
	}

	sub := s

	sub.min[0] += float32(region.Min.X)
	sub.min[1] += float32(region.Min.Y)

	size := region.Size()
	sub.size[0] = float32(size.X)
	sub.size[1] = float32(size.Y)

	return sub
}

// func (s Sprite) DrawRegion(region geo.Rectangle, transform geo.Transform, depth float32) {
// 	s.sheet.initialize()
// 	src := s.sheet.sources[s.index]
// 	dst := transform.Times(geo.Translation(src.offset))
// 	size := region.Size()
// 	s.sheet.renderer.DrawInstance(
// 		float32(dst[0][0]), float32(dst[0][1]), float32(dst[0][2]), float32(dst[1][0]), float32(dst[1][1]), float32(dst[1][2]),
// 		depth,
// 		src.location[0]+float32(region.Min.X-src.offset.X), src.location[1]+float32(region.Min.Y-src.offset.Y), src.location[2],
// 		float32(size.X), float32(size.Y),
// 	)
// }
