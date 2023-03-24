package image

import (
	"image"
	"image/color"
)

func MirrorX(i image.Image) image.Image {
	return mirroredXImage{i}
}

func MirrorY(i image.Image) image.Image {
	return mirroredYImage{i}
}

type mirroredXImage struct{ image.Image }

func (i mirroredXImage) Bounds() image.Rectangle {
	b := i.Image.Bounds()
	b.Max.X += b.Dx()
	return b
}

func (i mirroredXImage) At(x, y int) color.Color {
	max := i.Image.Bounds().Max.X
	if x > max {
		x = max - (x - max)
	}
	return i.Image.At(x, y)
}

type mirroredYImage struct{ image.Image }

func (i mirroredYImage) Bounds() image.Rectangle {
	b := i.Image.Bounds()
	b.Max.Y += b.Dy()
	return b
}

func (i mirroredYImage) At(x, y int) color.Color {
	max := i.Image.Bounds().Max.Y
	if y > max {
		y = max - (y - max)
	}
	return i.Image.At(x, y)
}
