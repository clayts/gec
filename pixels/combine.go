package pixels

import (
	"image"
	"image/color"
)

func CombineX(a, b image.Image) image.Image {
	if a.ColorModel() != b.ColorModel() {
		panic("images must have the same color model")
	}

	model := a.ColorModel()

	if a.Bounds().Min.Y != b.Bounds().Min.Y {
		panic("images must have the same Min.Y")
	}
	if a.Bounds().Max.Y != b.Bounds().Max.Y {
		panic("images must have the same Max.Y")
	}
	if a.Bounds().Min.X != b.Bounds().Min.X {
		panic("images must have the same Min.X")
	}
	bounds := a.Bounds()
	bounds.Max.X += b.Bounds().Dx()

	at := func(x, y int) color.Color {
		if x >= a.Bounds().Max.X {
			x -= a.Bounds().Max.X
			return b.At(x, y)
		}
		return a.At(x, y)
	}

	return NewProcedural(at, bounds, model)

}

func CombineY(a, b image.Image) image.Image {
	if a.ColorModel() != b.ColorModel() {
		panic("images must have the same color model")
	}

	model := a.ColorModel()

	if a.Bounds().Min.X != b.Bounds().Min.X {
		panic("images must have the same Min.X")
	}
	if a.Bounds().Max.X != b.Bounds().Max.X {
		panic("images must have the same Max.X")
	}
	if a.Bounds().Min.Y != b.Bounds().Min.Y {
		panic("images must have the same Min.Y")
	}
	bounds := a.Bounds()
	bounds.Max.Y += b.Bounds().Dy()

	at := func(x, y int) color.Color {
		if y >= a.Bounds().Max.Y {
			y -= a.Bounds().Max.Y
			return b.At(x, y)
		}
		return a.At(x, y)
	}

	return NewProcedural(at, bounds, model)

}
