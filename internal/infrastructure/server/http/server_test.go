package internalhttp

import (
	"encoding/json"
	"image"
	"image/color"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/alexandr-lakeev/otus-final-project/internal/app"
	"github.com/alexandr-lakeev/otus-final-project/internal/app/usecase"
	"github.com/alexandr-lakeev/otus-final-project/internal/config"
	internalcache "github.com/alexandr-lakeev/otus-final-project/internal/infrastructure/cache"
	internalimage "github.com/alexandr-lakeev/otus-final-project/internal/infrastructure/image"
	internallogger "github.com/alexandr-lakeev/otus-final-project/internal/infrastructure/logger"
	"github.com/stretchr/testify/require"
)

const TestHeader = "X-Extra-Header"

var headerValue string

func createServer() *http.Server {
	logger, err := internallogger.New(config.LoggerConf{Env: "test", Level: "INFO"})
	if err != nil {
		log.Fatal(err)
	}

	httpClient := &http.Client{
		Timeout: time.Second,
	}

	usecase := usecase.New(
		internalimage.NewLoader(httpClient),
		internalimage.NewResizer(),
		internalcache.NewCache(10, os.TempDir()),
		logger,
	)

	return NewServer(config.ServerConf{
		BindAddress: ":8080",
	}, usecase, logger)
}

func createFakeImageServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// store the proxied header value
		headerValue = r.Header.Get(TestHeader)

		if r.URL.Path == "/img/success/100x100" {
			err := jpeg.Encode(w, createTestImage(100, 100), &jpeg.Options{Quality: 100})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			return
		}

		if r.URL.Path == "/img/not-an-image" {
			response := map[string]string{
				"message": "this is not an image",
			}
			json.NewEncoder(w).Encode(response)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.URL.Path == "/img/error/400" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if r.URL.Path == "/img/error/500" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
}

// create image with pattern (f = #ffffff, 0 = #000000):
//
// f f f | 0 0 0
// f f f | 0 0 0
// f f f | 0 0 0
// ------x------
// 0 0 0 | f f f
// 0 0 0 | f f f
// 0 0 0 | f f f
//
func createTestImage(width, height int) *image.RGBA64 {
	img := image.NewRGBA64(image.Rect(0, 0, width, height))

	halfWidth := width / 2
	halfHeight := width / 2

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if (x > halfWidth && y <= halfHeight) || (x <= halfWidth && y > halfHeight) {
				img.Set(x, y, color.Black)
			} else {
				img.Set(x, y, color.White)
			}
		}
	}

	return img
}

func TestServer(t *testing.T) {
	t.Run("headers pass", func(t *testing.T) {
		imgServer := createFakeImageServer()
		defer imgServer.Close()

		imgServBaseUrl := url.QueryEscape(strings.Replace(imgServer.URL, "http://", "", 1))

		reqUrl := path.Join(
			"/fill/50/50",
			imgServBaseUrl,
			"/img/success/100x100",
		)

		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, reqUrl, nil)

		// this header must be proxied
		testHeaderValue := "test"
		req.Header.Set(TestHeader, testHeaderValue)

		createServer().Handler.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Result().StatusCode)

		// check the proxied header value
		require.Equal(t, testHeaderValue, headerValue)
	})

	t.Run("fill image", func(t *testing.T) {
		tests := []struct {
			name   string
			width  int
			height int
			imgUrl string
		}{
			{
				name:   "simple",
				width:  50,
				height: 50,
				imgUrl: "/img/success/100x100",
			},
			{
				name:   "only width",
				width:  50,
				height: 100,
				imgUrl: "/img/success/100x100",
			},
			{
				name:   "only height",
				width:  50,
				height: 100,
				imgUrl: "/img/success/100x100",
			},
			{
				name:   "smaller than fill size",
				width:  200,
				height: 200,
				imgUrl: "/img/success/100x100",
			},
		}

		for _, tc := range tests {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				imgServer := createFakeImageServer()
				defer imgServer.Close()

				imgServBaseUrl := url.QueryEscape(strings.Replace(imgServer.URL, "http://", "", 1))

				reqUrl := path.Join(
					"/fill",
					strconv.Itoa(tc.width),
					strconv.Itoa(tc.height),
					imgServBaseUrl,
					tc.imgUrl,
				)

				rec := httptest.NewRecorder()
				req, _ := http.NewRequest(http.MethodGet, reqUrl, nil)

				createServer().Handler.ServeHTTP(rec, req)

				require.Equal(t, http.StatusOK, rec.Result().StatusCode)

				// check image height and width
				img, _, err := image.Decode(rec.Body)
				require.NoError(t, err)

				bounds := img.Bounds()

				require.Equal(t, tc.width, bounds.Max.X)
				require.Equal(t, tc.height, bounds.Max.Y)

				// check image color in corners whith 5px padding
				padding := 5

				pixelColorTopLeft := color.Gray16Model.Convert(img.At(padding, padding))
				require.Equal(t, pixelColorTopLeft, color.White)

				pixelColorTopRight := color.Gray16Model.Convert(img.At(bounds.Max.X-padding, padding))
				require.Equal(t, pixelColorTopRight, color.Black)

				pixelColorBottomLeft := color.Gray16Model.Convert(img.At(padding, bounds.Max.Y-padding))
				require.Equal(t, pixelColorBottomLeft, color.Black)

				pixelColorBottomRight := color.Gray16Model.Convert(img.At(bounds.Max.X-padding, bounds.Max.Y-padding))
				require.Equal(t, pixelColorBottomRight, color.White)
			})
		}
	})

	t.Run("remote error", func(t *testing.T) {
		tests := []struct {
			name string
			url  string
			err  error
		}{
			{
				name: "not an image",
				url:  "/img/not-an-image",
				err:  app.ErrContentNotImage,
			},
			{
				name: "not found",
				url:  "/img/error/404",
				err:  app.ErrImageNotFound,
			},
			{
				name: "bad request",
				url:  "/img/error/400",
				err:  app.ErrBadRequest,
			},
			{
				name: "internal error",
				url:  "/img/error/500",
				err:  app.ErrInternal,
			},
		}

		var err error
		stdout := path.Join(os.TempDir(), "/stdout")

		os.Stdout, err = os.Create(stdout)
		if err != nil {
			log.Fatal(err)
		}

		imgServer := createFakeImageServer()
		defer imgServer.Close()

		imgServBaseUrl := url.QueryEscape(strings.Replace(imgServer.URL, "http://", "", 1))

		for _, tc := range tests {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				reqUrl := path.Join(
					"/fill/100/100",
					imgServBaseUrl,
					tc.url,
				)

				rec := httptest.NewRecorder()
				req, _ := http.NewRequest(http.MethodGet, reqUrl, nil)

				createServer().Handler.ServeHTTP(rec, req)

				require.Equal(t, http.StatusBadGateway, rec.Result().StatusCode)

				// check log message
				logContent, err := ioutil.ReadFile(stdout)
				require.NoError(t, err)
				require.Contains(t, string(logContent), tc.err.Error())
			})
		}
	})
}
