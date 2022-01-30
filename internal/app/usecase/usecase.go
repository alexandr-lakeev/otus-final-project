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
	img, ok := u.cache.Get(command.ImgUrl, command.Width, command.Height)
	if ok {
		u.logger.Info("got image from cache")
		return img, nil
	}

	img, err := u.loader.Load(ctx, command.ImgUrl, command.Headers)
	if err != nil {
		return nil, err
	}

	u.logger.Info("got image from remote")

	resizedImg := u.resizer.Fill(img, command.Width, command.Height)

	u.cache.Set(command.ImgUrl, command.Width, command.Height, resizedImg)

	return resizedImg, nil
}
