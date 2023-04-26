package pixels

import (
	"image"
	"image/color"
)

func FlipX(i image.Image) image.Image {
	model := i.ColorModel()

	bounds := i.Bounds()

	at := func(x, y int) color.Color {
		x = (i.Bounds().Max.X - 1) - x
		return i.At(x, y)
	}

	return NewProcedural(at, bounds, model)
}

func FlipY(i image.Image) image.Image {
	model := i.ColorModel()

	bounds := i.Bounds()

	at := func(x, y int) color.Color {
		y = (i.Bounds().Max.Y - 1) - y
		return i.At(x, y)
	}

	return NewProcedural(at, bounds, model)
}
