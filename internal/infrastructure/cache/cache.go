package internalcache

import (
	"container/list"
	"crypto/sha1"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/alexandr-lakeev/otus-final-project/internal/app"
)

type LruCache struct {
	capacity int
	queue    *list.List
	items    map[string]*list.Element
	dir      string
	lock     sync.Mutex
}

type CacheItem struct {
	Key    string
	Url    string
	Width  int
	Height int
	Path   string
}

func NewCache(capacity int, dir string) app.Cache {
	return &LruCache{
		capacity: capacity,
		queue:    list.New(),
		items:    make(map[string]*list.Element),
		dir:      dir,
	}
}

func (c *LruCache) Set(url string, width, height int, img image.Image) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	key := c.getKey(url, width, height)
	path, err := c.createPath(c.dir, key)
	if err != nil {
		return err
	}

	listItem, exists := c.items[key]

	if !exists {
		if c.queue.Len() == c.capacity {
			if err := c.delete(c.queue.Back()); err != nil {
				return err
			}
		}
	} else {
		c.queue.Remove(listItem)
	}

	err = c.saveToFile(path, img)
	if err != nil {
		return err
	}

	c.items[key] = c.queue.PushFront(&CacheItem{
		Key:    key,
		Url:    url,
		Width:  width,
		Height: height,
		Path:   path,
	})

	return nil
}

func (c *LruCache) Get(url string, width, height int) (image.Image, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	key := c.getKey(url, width, height)
	listItem, exists := c.items[key]

	if exists {
		path := listItem.Value.(*CacheItem).Path
		c.queue.MoveToFront(listItem)

		img, err := c.readFromFile(path)
		if err != nil {
			return nil, err
		}

		return img, nil
	}

	return nil, app.ErrNotFoundInCache
}

func (c *LruCache) delete(item *list.Element) error {
	c.queue.Remove(item)

	cacheItem := item.Value.(*CacheItem)
	delete(c.items, cacheItem.Key)

	return os.Remove(cacheItem.Path)
}

func (c *LruCache) getKey(url string, width, height int) string {
	return c.getHash(url + strconv.Itoa(width) + strconv.Itoa(height))
}

func (i *LruCache) getHash(key string) string {
	h := sha1.New()
	h.Write([]byte(key))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (i *LruCache) createPath(dir, key string) (string, error) {
	pathParts := append([]string{dir}, strings.Split(key, "")...)
	path := filepath.Join(pathParts...)

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return "", err
	}

	return filepath.Join(path, key), nil
}

func (c *LruCache) saveToFile(path string, img image.Image) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if err = jpeg.Encode(file, img, nil); err != nil {
		return err
	}

	return nil
}

func (c *LruCache) readFromFile(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	img, err := jpeg.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}
