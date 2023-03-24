package noise

import (
	"image"
	"image/color"
)

type ColorScheme struct {
	Model  color.Model
	Select func(v float64) color.Color
}

type Image struct {
	s2     Source2D
	bounds image.Rectangle
	scheme ColorScheme
}

func NewImage(s2 Source2D, bounds image.Rectangle, scheme ColorScheme) Image {
	return Image{s2: s2, bounds: bounds, scheme: scheme}
}

func (i Image) At(x, y int) color.Color {
	return i.scheme.Select(i.s2(float64(x), float64(y)))
}

func (i Image) Bounds() image.Rectangle {
	return i.bounds
}

func (i Image) ColorModel() color.Model {
	return i.scheme.Model
}
