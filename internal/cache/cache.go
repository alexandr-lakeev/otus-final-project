package cache

import (
	"container/list"
	"crypto/sha1"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
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

func (c *LruCache) Set(url string, width, height int, img image.Image) {
	c.lock.Lock()
	defer c.lock.Unlock()

	key := c.getKey(url, width, height)
	path, err := c.createPath(c.dir, key)
	// TODO resolve errors
	if err != nil {
		log.Println(err)
	}
	listItem, exists := c.items[key]

	if !exists {
		if c.queue.Len() == c.capacity {
			// TODO resolve errors
			if err := c.delete(c.queue.Back()); err != nil {
				log.Println(err)
			}
		}

		// TODO resolve errors
		err := c.saveToFile(path, img)
		if err != nil {
			log.Println(err)
		}
	} else {
		c.queue.Remove(listItem)
	}

	c.items[key] = c.queue.PushFront(&CacheItem{
		Key:    key,
		Url:    url,
		Width:  width,
		Height: height,
		Path:   path,
	})
}

func (c *LruCache) Get(url string, width, height int) (image.Image, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	key := c.getKey(url, width, height)
	listItem, exists := c.items[key]

	if exists {
		path := listItem.Value.(*CacheItem).Path
		c.queue.MoveToFront(listItem)

		img, err := c.readFromFile(path)
		// TODO resolve errors
		if err != nil {
			return nil, false
		}

		return img, true
	}

	return nil, false
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
	path := dir + "/" + strings.Join(strings.Split(key, ""), "/")
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return "", err
	}

	return path + "/" + key, nil
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
