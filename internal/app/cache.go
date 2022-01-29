package app

import "image"

type Cache interface {
	Get(url string) image.Image
	Set(url string, img image.Image)
}
