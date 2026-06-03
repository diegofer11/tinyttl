package tinyttl

import "time"

func (c *Cache) startCleanup() {
	ticker := time.NewTicker(c.cleanupInterval)
	defer ticker.Stop()
	defer close(c.doneChan)

	for {
		select {
		case <-ticker.C:
			c.deleteExpired()
		case <-c.stopChan:
			return
		}
	}
}

func (c *Cache) deleteExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key, item := range c.items {
		if c.isExpired(item) {
			delete(c.items, key)
		}
	}
}

func (c *Cache) Close() {
	if c.stopChan == nil || c.doneChan == nil {
		return
	}

	c.closeOnce.Do(func() {
		close(c.stopChan)
		<-c.doneChan
	})
}
