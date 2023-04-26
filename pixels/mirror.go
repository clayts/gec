package pixels

import (
	"image"
)

func MirrorX(i image.Image) image.Image {
	return CombineX(i, FlipX(i))
}

func MirrorY(i image.Image) image.Image {
	return CombineY(i, FlipY(i))
}
