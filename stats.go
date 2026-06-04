package tinyttl

type Stats struct {
	Sets        uint64
	Hits        uint64
	Misses      uint64
	Deletes     uint64
	Expirations uint64
}

type cacheStats struct {
	sets        uint64
	hits        uint64
	misses      uint64
	deletes     uint64
	expirations uint64
}

// Stats return a snapshot of the current cache statistics
func (c *Cache) Stats() Stats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return Stats{
		Sets:        c.stats.sets,
		Hits:        c.stats.hits,
		Misses:      c.stats.misses,
		Deletes:     c.stats.deletes,
		Expirations: c.stats.expirations,
	}
}
