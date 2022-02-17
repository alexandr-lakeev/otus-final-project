package internalimage

import (
	"bytes"
	"context"
	"image"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/alexandr-lakeev/otus-final-project/internal/app"
)

var statusCodeToError = map[int]error{
	http.StatusBadRequest:          app.ErrBadRequest,
	http.StatusNotFound:            app.ErrImageNotFound,
	http.StatusInternalServerError: app.ErrInternal,
	// TODO add more if needed
}

type ImageLoader struct {
	client *http.Client
}

func NewLoader(client *http.Client) *ImageLoader {
	return &ImageLoader{
		client: client,
	}
}

func (l *ImageLoader) Load(ctx context.Context, uri string, headers http.Header) (image.Image, error) {
	parsedUrl, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	parsedUrl.Scheme = "http"
	req, err := http.NewRequestWithContext(ctx, "GET", parsedUrl.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header = headers

	response, err := l.client.Do(req)
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
	if statusCode == http.StatusOK {
		return nil
	}

	err, ok := statusCodeToError[statusCode]
	if ok {
		return err
	}

	return app.ErrUnknown
}
