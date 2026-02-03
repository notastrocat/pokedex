package pokecache

import "time"

type cacheEntry struct {
	createdAt time.Time
	val []byte
}

type Cache struct {
	cache map[string]cacheEntry
}

func (c *Cache) reapLoop(updateInterval time.Duration) {
	ticker := time.NewTicker(updateInterval)
	go func() {
		for {
			<-ticker.C
			now := time.Now()
			for key, entry := range c.cache {
				if now.Sub(entry.createdAt) > updateInterval {
					delete(c.cache, key)
				}
			}
		}
	}()
}

func NewCache(updateInterval time.Duration) {
	cache.reapLoop(updateInterval)
}
