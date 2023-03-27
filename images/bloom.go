package images

import (
	"image"
	"image/color"

	"github.com/anthonynsimon/bild/blend"
	"github.com/anthonynsimon/bild/blur"
	"github.com/anthonynsimon/bild/effect"
)

func Bloom(img image.Image, radius int, glow float64) image.Image {

	imgSize := img.Bounds().Size()
	newSize := image.Rect(0, 0, imgSize.X+((radius+1)*2), imgSize.Y+((radius+1)*2))

	extended := image.NewRGBA(newSize)
	for x := 0; x < imgSize.X; x++ {
		for y := 0; y < imgSize.Y; y++ {
			extended.Set(radius+x+1, radius+y+1, img.At(x, y))
		}
	}

	newImg := effect.Dilate(extended, 1)

	newImg = blur.Gaussian(newImg, float64(radius))
	for x := 0; x < newSize.Max.X; x++ {
		for y := 0; y < newSize.Max.X; y++ {
			c := newImg.RGBAAt(x, y)
			c2 := color.NRGBA{c.R, c.G, c.B, uint8(255.0 * glow)}
			newImg.Set(x, y, c2)
		}
	}

	newImg = blend.Add(newImg, extended)

	return newImg
}
