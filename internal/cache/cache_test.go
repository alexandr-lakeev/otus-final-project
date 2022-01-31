package cache

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10, os.TempDir())

		// img := image.NewNRGBA(image.Rect(0, 0, 100, 100))

		_, ok := c.Get("www.img.ru/some-img.jpg", 100, 100)
		require.False(t, ok)
	})
}
