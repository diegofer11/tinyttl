package tinyttl

// Hooks defines optional callbacks triggered by cache events.
type Hooks struct {
	OnSet    func(key string, value any)
	OnDelete func(key string, value any)
	OnExpire func(key string, value any)
}
