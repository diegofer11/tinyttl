package tinyttl

import (
	"sync"
	"time"
)

type cacheItem struct {
	value     any
	expiresAt time.Time
}

// Cache is an in-memory TTL cache for string keys and arbitrary values.
//
// Expired items are removed lazily when accessed through Get, Has, or Len.
// If background cleanup is enabled, expired items are also removed periodically.
type Cache struct {
	mu              sync.RWMutex
	items           map[string]cacheItem
	stats           cacheStats
	hooks           Hooks
	cleanupInterval time.Duration
	stopChan        chan struct{}
	doneChan        chan struct{}
	closeOnce       sync.Once
}

// New creates a new Cache with the provided options.
func New(options ...Option) *Cache {
	cache := &Cache{
		items: make(map[string]cacheItem),
	}

	for _, option := range options {
		option(cache)
	}

	if cache.cleanupInterval > 0 {
		cache.stopChan = make(chan struct{})
		cache.doneChan = make(chan struct{})
		go cache.startCleanup()
	}

	return cache
}

// Set stores a value under the given key with the provided TTL.
//
// A positive TTL sets an expiration time for the item.
// A non-positive TTL means the item does not expire.
func (c *Cache) Set(key string, value any, ttl time.Duration) {
	var expiresAt time.Time
	if ttl > 0 {
		expiresAt = time.Now().Add(ttl)
	}

	c.mu.Lock()
	c.items[key] = cacheItem{
		value:     value,
		expiresAt: expiresAt,
	}
	c.stats.sets++
	onSet := c.hooks.OnSet
	c.mu.Unlock()

	if onSet != nil {
		onSet(key, value)
	}
}

// Get returns the value stored under the given key.
//
// It returns false if the key does not exist or if the item has expired.
// Expired items are removed lazily during the call.
func (c *Cache) Get(key string) (any, bool) {
	c.mu.Lock()

	item, exists := c.items[key]
	if !exists {
		c.stats.misses++
		c.mu.Unlock()
		return nil, false
	}

	if c.isExpired(item) {
		delete(c.items, key)
		c.stats.expirations++
		c.stats.misses++

		onExpire := c.hooks.OnExpire
		value := item.value
		c.mu.Unlock()

		if onExpire != nil {
			onExpire(key, value)
		}

		return nil, false
	}
	c.stats.hits++
	value := item.value
	c.mu.Unlock()

	return value, true
}

// Delete removes the value stored under the given key.
//
// If the key does not exist, Delete does nothing.
func (c *Cache) Delete(key string) {
	c.mu.Lock()

	item, exists := c.items[key]
	if !exists {
		c.mu.Unlock()
		return
	}

	delete(c.items, key)
	c.stats.deletes++
	onDelete := c.hooks.OnDelete

	value := item.value
	c.mu.Unlock()

	if onDelete != nil {
		onDelete(key, value)
	}
}

// Has reports whether a non-expired value exists for the given key.
//
// Expired items are removed lazily during the call.
func (c *Cache) Has(key string) bool {
	c.mu.Lock()

	item, exists := c.items[key]
	if !exists {
		c.mu.Unlock()
		return false
	}

	if c.isExpired(item) {
		delete(c.items, key)
		c.stats.expirations++
		onExpire := c.hooks.OnExpire
		value := item.value
		c.mu.Unlock()

		if onExpire != nil {
			onExpire(key, value)
		}
		return false
	}

	c.mu.Unlock()
	return true
}

// Len returns the number of non-expired items currently stored in the cache.
//
// Expired items are removed lazily during the call.
func (c *Cache) Len() int {
	c.mu.Lock()

	count := 0
	var expired []struct {
		key   string
		value any
	}

	for key, item := range c.items {
		if c.isExpired(item) {
			delete(c.items, key)
			c.stats.expirations++
			if c.hooks.OnExpire != nil {
				expired = append(expired, struct {
					key   string
					value any
				}{
					key:   key,
					value: item.value,
				})
			}
			continue
		}
		count++
	}

	onExpire := c.hooks.OnExpire
	c.mu.Unlock()

	if onExpire != nil {
		for _, item := range expired {
			onExpire(item.key, item.value)
		}
	}

	return count
}

func (c *Cache) isExpired(item cacheItem) bool {
	return !item.expiresAt.IsZero() && time.Now().After(item.expiresAt)
}
