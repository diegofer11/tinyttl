package tinyttl

// Hooks defines optional callbacks triggered by cache events.
type Hooks struct {
	// OnSet is called after a value is successfully stored in the cache.
	OnSet func(key string, value any)
	// OnDelete is called after a value is successfully removed from the cache.
	OnDelete func(key string, value any)
	// OnExpire is called after a value has expired and is removed from the cache.
	OnExpire func(key string, value any)
}
