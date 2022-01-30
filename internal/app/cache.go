package app

import "image"

type Cache interface {
	Get(url string, width, height int) (image.Image, bool)
	Set(url string, width, height int, img image.Image)
}
