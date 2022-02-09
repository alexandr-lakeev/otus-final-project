package internalcache

import (
	"image"
	"os"
	"testing"

	"github.com/alexandr-lakeev/otus-final-project/internal/app"
	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	img100x100 := image.NewNRGBA(image.Rect(0, 0, 100, 100))
	img200x200 := image.NewNRGBA(image.Rect(0, 0, 200, 200))
	img300x300 := image.NewNRGBA(image.Rect(0, 0, 300, 300))
	img400x400 := image.NewNRGBA(image.Rect(0, 0, 400, 400))
	img500x500 := image.NewNRGBA(image.Rect(0, 0, 500, 500))
	img600x600 := image.NewNRGBA(image.Rect(0, 0, 600, 600))

	errNotFound := app.ErrNotFoundInCache

	t.Run("empty cache", func(t *testing.T) {
		cache := NewCache(5, os.TempDir())

		_, err := cache.Get("www.img.ru/some-img.jpg", 100, 100)

		require.ErrorIs(t, err, errNotFound)
	})

	t.Run("simple caching", func(t *testing.T) {
		cache := NewCache(5, os.TempDir())

		err := cache.Set("www.img.ru/some-img.jpg", 100, 100, img100x100)
		require.NoError(t, err)

		err = cache.Set("www.img.ru/some-img.jpg", 200, 200, img200x200)
		require.NoError(t, err)

		img100x100Cached, err := cache.Get("www.img.ru/some-img.jpg", 100, 100)

		require.NoError(t, err)
		require.Equal(t, 100, img100x100Cached.Bounds().Max.X)
		require.Equal(t, 100, img100x100Cached.Bounds().Max.Y)

		img200x200Cached, err := cache.Get("www.img.ru/some-img.jpg", 200, 200)

		require.NoError(t, err)
		require.Equal(t, 200, img200x200Cached.Bounds().Max.X)
		require.Equal(t, 200, img200x200Cached.Bounds().Max.Y)

		_, err = cache.Get("www.img.ru/some-img.jpg", 300, 300)

		require.ErrorIs(t, err, errNotFound)
	})

	t.Run("first added is removing", func(t *testing.T) {
		cache := NewCache(5, os.TempDir())

		err := cache.Set("www.img.ru/some-img.jpg", 100, 100, img100x100)
		require.NoError(t, err)

		err = cache.Set("www.img.ru/some-img.jpg", 200, 200, img200x200)
		require.NoError(t, err)

		err = cache.Set("www.img.ru/some-img.jpg", 300, 300, img300x300)
		require.NoError(t, err)

		err = cache.Set("www.img.ru/some-img.jpg", 400, 400, img400x400)
		require.NoError(t, err)

		err = cache.Set("www.img.ru/some-img.jpg", 500, 500, img500x500)
		require.NoError(t, err)

		err = cache.Set("www.img.ru/some-img.jpg", 600, 600, img600x600)
		require.NoError(t, err)

		img600x600Cached, err := cache.Get("www.img.ru/some-img.jpg", 600, 600)

		require.NoError(t, err)
		require.Equal(t, 600, img600x600Cached.Bounds().Max.X)
		require.Equal(t, 600, img600x600Cached.Bounds().Max.Y)

		_, err = cache.Get("www.img.ru/some-img.jpg", 100, 100)

		require.ErrorIs(t, err, errNotFound)
	})

	t.Run("first touched is removing", func(t *testing.T) {
		cache := NewCache(5, os.TempDir())

		err := cache.Set("www.img.ru/some-img.jpg", 100, 100, img100x100)
		require.NoError(t, err)

		err = cache.Set("www.img.ru/some-img.jpg", 200, 200, img200x200)
		require.NoError(t, err)

		err = cache.Set("www.img.ru/some-img.jpg", 300, 300, img300x300)
		require.NoError(t, err)

		err = cache.Set("www.img.ru/some-img.jpg", 400, 400, img400x400)
		require.NoError(t, err)

		err = cache.Set("www.img.ru/some-img.jpg", 500, 500, img500x500)
		require.NoError(t, err)

		_, err = cache.Get("www.img.ru/some-img.jpg", 500, 500)
		require.NoError(t, err)

		_, err = cache.Get("www.img.ru/some-img.jpg", 400, 400)
		require.NoError(t, err)

		_, err = cache.Get("www.img.ru/some-img.jpg", 300, 300)
		require.NoError(t, err)

		_, err = cache.Get("www.img.ru/some-img.jpg", 200, 200)
		require.NoError(t, err)

		_, err = cache.Get("www.img.ru/some-img.jpg", 100, 100)
		require.NoError(t, err)

		err = cache.Set("www.img.ru/some-img.jpg", 600, 600, img600x600)
		require.NoError(t, err)

		_, err = cache.Get("www.img.ru/some-img.jpg", 600, 600)
		require.NoError(t, err)

		_, err = cache.Get("www.img.ru/some-img.jpg", 500, 500)

		require.ErrorIs(t, err, errNotFound)
	})
}
