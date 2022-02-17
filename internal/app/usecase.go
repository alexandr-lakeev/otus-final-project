package app

import (
	"context"
	"image"
	"net/http"
)

type UseCase interface {
	Fill(context.Context, *FillCommand) (image.Image, error)
}

type FillCommand struct {
	ImgUrl  string
	Width   int
	Height  int
	Headers http.Header
}
