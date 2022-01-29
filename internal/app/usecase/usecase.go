package usecase

import (
	"context"
	"image"

	"github.com/alexandr-lakeev/otus-final-project/internal/app"
)

type UseCase struct {
	loader app.ImageLoader
	cache  app.Cache
	logger app.Logger
}

func New(loader app.ImageLoader, cache app.Cache, logger app.Logger) *UseCase {
	return &UseCase{
		loader: loader,
		cache:  cache,
		logger: logger,
	}
}

func (u *UseCase) Fill(ctx context.Context, command *app.FillCommand) (image.Image, error) {
	// TODO resize, cache
	return u.loader.Load(ctx, command.ImgUrl, command.Headers)
}
