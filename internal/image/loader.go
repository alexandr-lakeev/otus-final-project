package image

import (
	"bytes"
	"context"
	"image"
	"io"
	"net/http"
	"time"
)

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

	if response.StatusCode != 200 {
		return nil, err
	}

	imgBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(bytes.NewReader(imgBytes))

	return img, err
}
