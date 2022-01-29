package deliveryhttp

import (
	"context"
	"image/jpeg"
	"net/http"
	"strconv"
	"strings"

	"github.com/alexandr-lakeev/otus-final-project/internal/app"
)

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
		parts := strings.SplitN(r.URL.Path, "/", 5)

		if len(parts) < 5 {
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
			ImgUrl:  url,
			Width:   width,
			Height:  height,
			Headers: r.Header,
		})

		if err != nil {
			h.logger.Error("handler: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// TODO png
		w.Header().Set("Content-Type", "image/jpg")
		w.WriteHeader(http.StatusOK)

		err = jpeg.Encode(w, image, &jpeg.Options{
			Quality: 100,
		})

		if err != nil {
			h.logger.Error("handler: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
