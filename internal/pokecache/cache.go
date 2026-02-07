package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	cache map[string]cacheEntry
	mu    sync.RWMutex
	stop  chan struct{}
}

func (c *Cache) reapLoop(updateInterval time.Duration) {
	ticker := time.NewTicker(updateInterval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				c.mu.Lock()
				now := time.Now()
				for key, entry := range c.cache {
					if now.Sub(entry.createdAt) > updateInterval {
						delete(c.cache, key)
					}
				}
				c.mu.Unlock()	// we do NOT use defer in a loop!
			case <-c.stop:
				return
			}
		}
	}()
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, exists := c.cache[key]
	if !exists {
		return nil, false
	}
	return entry.val, true
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Close() {
	close(c.stop)
}

func NewCache(updateInterval time.Duration) *Cache {
	cache := &Cache{
		cache: make(map[string]cacheEntry),
		stop:  make(chan struct{}),
	}
	cache.reapLoop(updateInterval)
	return cache
}
