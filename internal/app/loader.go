package app

import (
	"context"
	"errors"
	"image"
	"net/http"
)

var ErrImageNotFound = errors.New("image not found")
var ErrBadRequest = errors.New("bad request")
var ErrInternal = errors.New("an internal error occurred while loading image")
var ErrUnknown = errors.New("an unknown error occurred while loading image")
var ErrContentNotImage = errors.New("content not an image")

type ImageLoader interface {
	Load(ctx context.Context, url string, headers http.Header) (image.Image, error)
}
