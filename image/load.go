package image

import (
	"fmt"
	"image"
	"image/draw"
	"os"
)

func LoadRGBA(file string) *image.RGBA {
	imgFile, err := os.Open(file)
	if err != nil {
		panic(fmt.Errorf("texture %q not found on disk: %v", file, err))
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		panic(err)
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		panic("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	return rgba
}
