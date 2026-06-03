package tinyttl

import "time"

type Option func(*Cache)

func WithCleanupInterval(interval time.Duration) Option {
	return func(c *Cache) {
		if interval > 0 {
			c.cleanupInterval = interval
		}
	}
}
