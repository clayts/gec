package images

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
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

func SaveImage(filename string, img image.Image) {
	imgFile, err := os.Create(filename)
	defer func() {
		if err := imgFile.Close(); err != nil {
			panic(err)
		}
	}()
	if err != nil {
		panic(err)
	}
	png.Encode(imgFile, img)
}
