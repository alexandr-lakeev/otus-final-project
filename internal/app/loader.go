package app

import (
	"context"
	"image"
	"net/http"
)

type ImageLoader interface {
	Load(ctx context.Context, url string, headers http.Header) (image.Image, error)
}
