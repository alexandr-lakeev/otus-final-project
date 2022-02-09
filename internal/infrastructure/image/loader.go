package internalimage

import (
	"bytes"
	"context"
	"image"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/alexandr-lakeev/otus-final-project/internal/app"
)

var statusCodeToError = map[int]error{
	400: app.ErrBadRequest,
	404: app.ErrImageNotFound,
	500: app.ErrInternal,
	// TODO add more if needed
}

type ImageLoader struct {
}

func NewLoader() *ImageLoader {
	return &ImageLoader{}
}

func (l *ImageLoader) Load(ctx context.Context, url string, headers http.Header) (image.Image, error) {
	req, err := http.NewRequest("GET", "http://"+url, nil)
	if err != nil {
		return nil, err
	}
	req.Header = headers

	client := http.Client{
		Timeout: 1 * time.Second,
	}

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if err := l.resolveStatusCode(response.StatusCode); err != nil {
		return nil, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if !l.isImage(body) {
		return nil, app.ErrContentNotImage
	}

	img, _, err := image.Decode(bytes.NewReader(body))

	return img, err
}

func (l *ImageLoader) isImage(body []byte) bool {
	return strings.Split(http.DetectContentType(body), "/")[0] == "image"
}

func (l *ImageLoader) resolveStatusCode(statusCode int) error {
	if statusCode == 200 {
		return nil
	}

	err, ok := statusCodeToError[statusCode]
	if ok {
		return err
	}

	return app.ErrUnknown
}
