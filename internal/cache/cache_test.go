package cache

import (
	"image"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	img100x100 := image.NewNRGBA(image.Rect(0, 0, 100, 100))
	img200x200 := image.NewNRGBA(image.Rect(0, 0, 200, 200))
	img300x300 := image.NewNRGBA(image.Rect(0, 0, 300, 300))
	img400x400 := image.NewNRGBA(image.Rect(0, 0, 400, 400))
	img500x500 := image.NewNRGBA(image.Rect(0, 0, 500, 500))
	img600x600 := image.NewNRGBA(image.Rect(0, 0, 600, 600))

	t.Run("empty cache", func(t *testing.T) {
		cache := NewCache(5, os.TempDir())

		_, ok := cache.Get("www.img.ru/some-img.jpg", 100, 100)

		require.False(t, ok)
	})

	t.Run("simple caching", func(t *testing.T) {
		cache := NewCache(5, os.TempDir())

		cache.Set("www.img.ru/some-img.jpg", 100, 100, img100x100)
		cache.Set("www.img.ru/some-img.jpg", 200, 200, img200x200)

		img100x100Cached, ok := cache.Get("www.img.ru/some-img.jpg", 100, 100)

		require.True(t, ok)
		require.Equal(t, 100, img100x100Cached.Bounds().Max.X)
		require.Equal(t, 100, img100x100Cached.Bounds().Max.Y)

		img200x200Cached, ok := cache.Get("www.img.ru/some-img.jpg", 200, 200)

		require.True(t, ok)
		require.Equal(t, 200, img200x200Cached.Bounds().Max.X)
		require.Equal(t, 200, img200x200Cached.Bounds().Max.Y)

		_, ok = cache.Get("www.img.ru/some-img.jpg", 300, 300)

		require.False(t, ok)
	})

	t.Run("last added is removing", func(t *testing.T) {
		cache := NewCache(5, os.TempDir())

		cache.Set("www.img.ru/some-img.jpg", 100, 100, img100x100)
		cache.Set("www.img.ru/some-img.jpg", 200, 200, img200x200)
		cache.Set("www.img.ru/some-img.jpg", 300, 300, img300x300)
		cache.Set("www.img.ru/some-img.jpg", 400, 400, img400x400)
		cache.Set("www.img.ru/some-img.jpg", 500, 500, img500x500)
		cache.Set("www.img.ru/some-img.jpg", 600, 600, img600x600)

		img600x600Cached, ok := cache.Get("www.img.ru/some-img.jpg", 600, 600)

		require.True(t, ok)
		require.Equal(t, 600, img600x600Cached.Bounds().Max.X)
		require.Equal(t, 600, img600x600Cached.Bounds().Max.Y)

		_, ok = cache.Get("www.img.ru/some-img.jpg", 100, 100)

		require.False(t, ok)
	})

	t.Run("first touched is removing", func(t *testing.T) {
		cache := NewCache(5, os.TempDir())

		cache.Set("www.img.ru/some-img.jpg", 100, 100, img100x100)
		cache.Set("www.img.ru/some-img.jpg", 200, 200, img200x200)
		cache.Set("www.img.ru/some-img.jpg", 300, 300, img300x300)
		cache.Set("www.img.ru/some-img.jpg", 400, 400, img400x400)
		cache.Set("www.img.ru/some-img.jpg", 500, 500, img500x500)

		_, ok := cache.Get("www.img.ru/some-img.jpg", 500, 500)
		require.True(t, ok)

		_, ok = cache.Get("www.img.ru/some-img.jpg", 400, 400)
		require.True(t, ok)

		_, ok = cache.Get("www.img.ru/some-img.jpg", 300, 300)
		require.True(t, ok)

		_, ok = cache.Get("www.img.ru/some-img.jpg", 200, 200)
		require.True(t, ok)

		_, ok = cache.Get("www.img.ru/some-img.jpg", 100, 100)
		require.True(t, ok)

		cache.Set("www.img.ru/some-img.jpg", 600, 600, img600x600)
		_, ok = cache.Get("www.img.ru/some-img.jpg", 600, 600)
		require.True(t, ok)

		_, ok = cache.Get("www.img.ru/some-img.jpg", 500, 500)

		require.False(t, ok)
	})
}
