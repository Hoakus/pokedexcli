package pokecache

import (
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	entry map[string]cacheEntry
	mu    sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	newCache := &Cache{
		entry: make(map[string]cacheEntry),
		mu:    sync.Mutex{},
	}
	go newCache.reapLoop(interval)
	return newCache
}

func (c *Cache) reapLoop(interval time.Duration) {
	for ; ; time.Tick(interval) {
		c.mu.Lock()
		for key, currentEntry := range c.entry {
			if time.Since(currentEntry.createdAt) > interval {
				delete(c.entry, key)
			}
		}
		c.mu.Unlock()
	}

}

func (c *Cache) Add(key string, value []byte) {
	if key == "" {
		fmt.Println("error adding to cache: no key")
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	newCacheEntry := cacheEntry{createdAt: time.Now(), val: value}
	c.entry[key] = newCacheEntry

	return
}

func (c *Cache) Get(key string) ([]byte, bool) {
	if key == "" {
		fmt.Println("error getting from cache: no key")
		return []byte{}, false
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if entry, ok := c.entry[key]; !ok {
		return []byte{}, false
	} else {
		return entry.val, true
	}
}
