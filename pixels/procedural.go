package pixels

import (
	"image"
	"image/color"
)

func NewProcedural(at func(x, y int) color.Color, bounds image.Rectangle, model color.Model) image.Image {
	return procedural{at: at, bounds: bounds, model: model}
}

func NewUniform(c color.Color, bounds image.Rectangle, model color.Model) image.Image {
	return NewProcedural(func(x, y int) color.Color { return c }, bounds, model)
}

type procedural struct {
	at     func(x, y int) color.Color
	bounds image.Rectangle
	model  color.Model
}

func (p procedural) At(x, y int) color.Color {
	return p.at(x, y)
}

func (p procedural) Bounds() image.Rectangle {
	return p.bounds
}

func (p procedural) ColorModel() color.Model {
	return p.model
}
