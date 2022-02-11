package deliveryhttp

import (
	"context"
	"image/jpeg"
	"net/http"
	"strconv"
	"strings"

	"github.com/alexandr-lakeev/otus-final-project/internal/app"
	"github.com/pkg/errors"
)

const UrlPartsQuantityBeforeImgPath = 5
const ImageQualityPercent = 100

type Handler struct {
	useCase app.UseCase
	logger  app.Logger
}

func NewHandler(useCase app.UseCase, logger app.Logger) *Handler {
	return &Handler{
		useCase: useCase,
		logger:  logger,
	}
}

func (h *Handler) Fill(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parts := strings.SplitN(r.URL.Path, "/", UrlPartsQuantityBeforeImgPath)

		if len(parts) < UrlPartsQuantityBeforeImgPath {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		width, err := strconv.Atoi(parts[2])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		height, err := strconv.Atoi(parts[3])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		url := parts[4]
		if url == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		image, err := h.useCase.Fill(ctx, &app.FillCommand{
			ImgUrl:  "//" + url, // to prevent error if target is ip address + port https://github.com/golang/go/issues/19297#issuecomment-282650053
			Width:   width,
			Height:  height,
			Headers: r.Header,
		})

		if err != nil {
			h.logger.Error(errors.Wrap(err, "fill error").Error())
			w.WriteHeader(http.StatusBadGateway)
			return
		}

		w.Header().Set("Content-Type", "image/jpg")
		w.WriteHeader(http.StatusOK)

		err = jpeg.Encode(w, image, &jpeg.Options{
			Quality: ImageQualityPercent,
		})

		if err != nil {
			h.logger.Error(errors.Wrap(err, "jpeg encode error").Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
