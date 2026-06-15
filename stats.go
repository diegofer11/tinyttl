package tinyttl

// Stats contain a snapshot of cache usage statistics.
type Stats struct {
	// Sets is the number of Set calls that stored a value in the cache.
	Sets uint64
	// Hits is the number of successful Get calls.
	Hits uint64
	// Misses is the number of Get calls that did not find a value in the cache.
	Misses uint64
	// Deletes is the number of Delete calls that removed a value from the cache.
	Deletes uint64
	// Expirations is the number of expired items that were removed from the cache.
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
