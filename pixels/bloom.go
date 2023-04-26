package pixels

import (
	"image"
	"image/color"

	"github.com/anthonynsimon/bild/blend"
	"github.com/anthonynsimon/bild/blur"
	"github.com/anthonynsimon/bild/effect"
)

func Bloom(img image.Image, radius int, glow float64) image.Image {

	newRect := image.Rect(img.Bounds().Min.X-radius, img.Bounds().Min.Y-radius, img.Bounds().Max.X+radius, img.Bounds().Max.Y+radius)

	extended := image.NewRGBA(newRect)
	for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
		for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
			extended.Set(x, y, img.At(x, y))
		}
	}

	newImg := effect.Dilate(extended, 2)

	newImg = blur.Gaussian(newImg, float64(radius))
	for x := newRect.Min.X; x < newRect.Max.X; x++ {
		for y := newRect.Min.Y; y < newRect.Max.Y; y++ {
			c := newImg.RGBAAt(x, y)
			c2 := color.NRGBA{c.R, c.G, c.B, uint8(255.0 * glow)}
			newImg.Set(x, y, c2)
		}
	}

	newImg = blend.Add(newImg, extended)

	return newImg
}
