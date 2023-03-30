package sprites

import (
	"image"

	geo "github.com/clayts/gec/geometry"
)

type Sprite struct {
	renderer *Sheet
	index    int
}

func (r *Sheet) MakeSprite(img image.Image) Sprite {
	r.sources = append(r.sources, struct {
		location [3]float32
		size     [2]float32
		image    image.Image
	}{
		image: img,
		size:  [2]float32{float32(img.Bounds().Dx()), float32(img.Bounds().Dy())},
	})
	return Sprite{
		renderer: r,
		index:    len(r.sources) - 1,
	}
}

func (s Sprite) Draw(dst geo.Transform, depth float32) {
	s.renderer.initialize()
	src := s.renderer.sources[s.index]
	s.renderer.renderer.DrawInstance(
		float32(dst[0][0]), float32(dst[0][1]), float32(dst[0][2]), float32(dst[1][0]), float32(dst[1][1]), float32(dst[1][2]),
		depth,
		src.location[0], src.location[1], src.location[2],
		src.size[0], src.size[1],
	)
}

func (s Sprite) Size() geo.Vector {
	src := s.renderer.sources[s.index]
	return geo.V(float64(src.size[0]), float64(src.size[1]))
}

func (s Sprite) Bounds() geo.Rectangle {
	return geo.R(geo.V(0, 0), s.Size())
}
