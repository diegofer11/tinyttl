package tinyttl

import "time"

// Option configures a Cache during initialization.
type Option func(*Cache)

// WithCleanupInterval enables background cleanup of expired items.
//
// If the interval is greater than zero, the cache starts a background goroutine
// that periodically removes expired items.
// If the interval is zero or negative, background cleanup is disabled.
func WithCleanupInterval(interval time.Duration) Option {
	return func(c *Cache) {
		if interval > 0 {
			c.cleanupInterval = interval
		}
	}
}

// WithHooks configures hooks for cache events.
func WithHooks(hooks Hooks) Option {
	return func(c *Cache) {
		c.hooks = hooks
	}
}
