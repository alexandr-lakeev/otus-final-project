package internalhttp

import (
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/alexandr-lakeev/otus-final-project/internal/app"
	"github.com/alexandr-lakeev/otus-final-project/internal/app/usecase"
	"github.com/alexandr-lakeev/otus-final-project/internal/config"
	internalcache "github.com/alexandr-lakeev/otus-final-project/internal/infrastructure/cache"
	internalimage "github.com/alexandr-lakeev/otus-final-project/internal/infrastructure/image"
	internallogger "github.com/alexandr-lakeev/otus-final-project/internal/infrastructure/logger"
	"github.com/stretchr/testify/require"
)

func createServer() *Server {
	logger, err := internallogger.New(config.LoggerConf{Env: "test", Level: "INFO"})
	if err != nil {
		log.Fatal(err)
	}

	usecase := usecase.New(
		internalimage.NewLoader(),
		internalimage.NewResizer(),
		internalcache.NewCache(10, os.TempDir()),
		logger,
	)

	return NewServer(config.ServerConf{
		BindAddress: ":8080",
	}, usecase, logger)
}

func createImageFakeServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/img/success/200x200" {
			err := jpeg.Encode(w, image.NewRGBA(image.Rect(0, 0, 200, 200)), &jpeg.Options{Quality: 100})
			if err != nil {
				w.WriteHeader(500)
			}
			w.WriteHeader(200)
			return
		}

		w.WriteHeader(404)
	}))
}

func TestServer(t *testing.T) {
	t.Run("fill image", func(t *testing.T) {
		imgServer := createImageFakeServer()
		defer imgServer.Close()

		imageReqUrl := url.QueryEscape(strings.Replace(imgServer.URL, "http://", "", 1)) + "/img/success/200x200"

		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/fill/100/100/"+imageReqUrl, nil)

		createServer().ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Result().StatusCode)

		// check image height and width
		img, _, err := image.Decode(rec.Body)
		require.NoError(t, err)

		bounds := img.Bounds()

		require.Equal(t, 100, bounds.Max.X)
		require.Equal(t, 100, bounds.Max.Y)
	})

	t.Run("image not found", func(t *testing.T) {
		var err error
		stdout := os.TempDir() + "/stdout"

		os.Stdout, err = os.Create(stdout)
		if err != nil {
			log.Fatal(err)
		}

		imgServer := createImageFakeServer()
		defer imgServer.Close()

		imageReqUrl := url.QueryEscape(strings.Replace(imgServer.URL, "http://", "", 1)) + "/img/not-found"

		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/fill/100/100/"+imageReqUrl, nil)

		createServer().ServeHTTP(rec, req)

		require.Equal(t, http.StatusBadGateway, rec.Result().StatusCode)

		// check log message - image not found
		errorMessage := app.ErrImageNotFound.Error()
		logContent, err := ioutil.ReadFile(stdout)
		require.NoError(t, err)

		require.Contains(t, string(logContent), errorMessage)
	})
}
