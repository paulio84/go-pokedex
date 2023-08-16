package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	entries map[string]cacheEntry
	mu      *sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		entries: make(map[string]cacheEntry),
		mu:      &sync.Mutex{},
	}

	if interval > 0 {
		go c.reapLoop(interval)
	}

	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) (response []byte, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	response = c.entries[key].val
	if response == nil {
		return nil, false
	}

	return response, true
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)

	for range ticker.C {
		c.deleteExpired(interval)
	}

	// for {
	// 	select {
	// 	case <-ticker.C:
	// 		c.deleteExpired(interval)
	// 	}
	// }
}

func (c *Cache) deleteExpired(interval time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key, entry := range c.entries {
		diff := time.Since(entry.createdAt)
		if diff >= interval {
			delete(c.entries, key)
		}
	}
}
