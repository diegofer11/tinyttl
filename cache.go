package tinyttl

import (
	"sync"
	"time"
)

type cacheItem struct {
	value     any
	expiresAt time.Time
}

type Cache struct {
	mu    sync.RWMutex
	items map[string]cacheItem
}

func New() *Cache {
	return &Cache{
		items: make(map[string]cacheItem),
	}
}

func (c *Cache) Set(key string, value any, ttl time.Duration) {
	var expiresAt time.Time
	if ttl > 0 {
		expiresAt = time.Now().Add(ttl)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = cacheItem{
		value:     value,
		expiresAt: expiresAt,
	}
}

func (c *Cache) Get(key string) (any, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, exists := c.items[key]
	if !exists {
		return nil, false
	}

	if c.isExpired(item) {
		delete(c.items, key)
		return nil, false
	}

	return item.value, true
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

func (c *Cache) Has(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, exists := c.items[key]
	if !exists {
		return false
	}

	if c.isExpired(item) {
		delete(c.items, key)
		return false
	}

	return true
}

func (c *Cache) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	count := 0
	for key, item := range c.items {
		if c.isExpired(item) {
			delete(c.items, key)
			continue
		}
		count++
	}
	return count
}

func (c *Cache) isExpired(item cacheItem) bool {
	return !item.expiresAt.IsZero() && time.Now().After(item.expiresAt)
}
