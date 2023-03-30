package procedural

import (
	"image"
	"image/color"
)

type Image struct {
	get    func(x, y int) color.Color
	bounds image.Rectangle
	model  color.Model
}

func NewImage(get func(x, y int) color.Color, bounds image.Rectangle, model color.Model) Image {
	return Image{get: get, bounds: bounds, model: model}
}

func (i Image) At(x, y int) color.Color {
	return i.get(x, y)
}

func (i Image) Bounds() image.Rectangle {
	return i.bounds
}

func (i Image) ColorModel() color.Model {
	return i.model
}

func NewUniformImage(c color.Color, bounds image.Rectangle, model color.Model) Image {
	return NewImage(func(x, y int) color.Color { return c }, bounds, model)
}
