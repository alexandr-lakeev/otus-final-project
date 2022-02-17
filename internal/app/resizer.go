package app

import "image"

type ImageResizer interface {
	Fill(img image.Image, width, height int) image.Image
}
