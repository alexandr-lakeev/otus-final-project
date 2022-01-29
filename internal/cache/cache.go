package cache

import "image"

type Cache struct {
}

func New() *Cache {
	return &Cache{}
}

func (c *Cache) Get(url string) image.Image {
	return nil
}

func (c *Cache) Set(url string, img image.Image) {

}
