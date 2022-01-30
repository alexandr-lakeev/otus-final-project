package cache

import (
	"testing"

	drivermemory "github.com/alexandr-lakeev/otus-final-project/internal/cache/driver/memory"
	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10, drivermemory.New())

		// img := image.NewNRGBA(image.Rect(0, 0, 100, 100))

		_, ok := c.Get("www.img.ru/some-img.jpg", 100, 100)
		require.False(t, ok)
	})
}
