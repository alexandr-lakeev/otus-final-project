package app

import (
	"errors"
	"image"
)

var ErrNotFoundInCache = errors.New("not found in cache")

type Cache interface {
	Get(url string, width, height int) (image.Image, error)
	Set(url string, width, height int, img image.Image) error
}
