package usecase

import (
	"context"
	"image"

	"github.com/alexandr-lakeev/otus-final-project/internal/app"
)

type UseCase struct {
	loader  app.ImageLoader
	resizer app.ImageResizer
	cache   app.Cache
	logger  app.Logger
}

func New(loader app.ImageLoader, resizer app.ImageResizer, cache app.Cache, logger app.Logger) *UseCase {
	return &UseCase{
		loader:  loader,
		resizer: resizer,
		cache:   cache,
		logger:  logger,
	}
}

func (u *UseCase) Fill(ctx context.Context, command *app.FillCommand) (image.Image, error) {
	// TODO cache
	img, err := u.loader.Load(ctx, command.ImgUrl, command.Headers)
	if err != nil {
		return nil, err
	}

	return u.resizer.Fill(img, command.Width, command.Height), nil
}
