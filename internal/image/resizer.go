package image

import (
	"image"

	"github.com/disintegration/imaging"
)

type ImageResizer struct {
}

func NewResizer() *ImageResizer {
	return &ImageResizer{}
}

func (r *ImageResizer) Fill(img image.Image, width, height int) image.Image {
	return imaging.Fill(img, width, height, imaging.Center, imaging.Lanczos)
}
