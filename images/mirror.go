package images

import (
	"image"
	"image/color"
)

func MirrorX(i image.Image) image.Image {
	model := i.ColorModel()

	bounds := i.Bounds()
	bounds.Max.X += bounds.Dx()

	at := func(x, y int) color.Color {
		max := i.Bounds().Max.X
		if x >= max {
			x = (max - (x - max)) - 1
		}
		return i.At(x, y)
	}

	return NewProcedural(at, bounds, model)
}

func MirrorY(i image.Image) image.Image {
	model := i.ColorModel()

	bounds := i.Bounds()
	bounds.Max.Y += bounds.Dy()

	at := func(x, y int) color.Color {
		max := i.Bounds().Max.Y
		if y >= max {
			y = (max - (y - max)) - 1
		}
		return i.At(x, y)
	}

	return NewProcedural(at, bounds, model)
}
