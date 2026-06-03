# tinyttl

`tinyttl` is a lightweight in-memory TTL cache for Go.

It provides a simple API for storing values with optional expiration times, lazy expiration on access, and optional background cleanup for expired items.

## Features

- In-memory cache with string keys
- Store arbitrary values
- Per-item TTL support
- Lazy expiration on `Get`, `Has`, and `Len`
- Optional background cleanup
- Simple and minimal API

## Installation

```bash
go get github.com/diegofer11/tinyttl
```

## Quick start

```go
package main

import (
	"fmt"
	"time"

	"github.com/diegofer11/tinyttl"
)

func main() {
	cache := tinyttl.New()

	cache.Set("message", "hello", time.Minute)

	value, found := cache.Get("message")
	if !found {
		fmt.Println("value not found")
		return
	}

	fmt.Println(value)
}
```

## Expiration behavior

A positive TTL sets an expiration time for the item.

A non-positive TTL means the item does not expire.

Expired items are removed lazily when accessed through:
- `Get`
- `Has`
- `Len`

## Background cleanup

You can enable periodic cleanup of expired items with `WithCleanupInterval`:

```go
package main

import (
	"time"

	"github.com/diegofer11/tinyttl"
)

func main() {
	cache := tinyttl.New(tinyttl.WithCleanupInterval(30 * time.Second))
	defer cache.Close()

	cache.Set("session", "abc123", time.Minute)
}
```

If background cleanup is enabled, remember to call `Close()` when the cache is no longer needed.

## API overview

### Create a cache

```go
cache := tinyttl.New()
```

### Store a value

```go
cache.Set("user:1", "Diego", time.Minute)
```

### Store a value without expiration

```go
cache.Set("config", "enabled", 0)
```

### Read a value

```go
value, found := cache.Get("user:1")
if found {
	fmt.Println(value)
}
```

### Check whether a key exists

```go
if cache.Has("user:1") {
	fmt.Println("key exists")
}
```

### Delete a key

```go
cache.Delete("user:1")
```

### Count active items

```go
fmt.Println(cache.Len())
```

## Development

Run tests:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

Open the HTML coverage report:

```bash
go tool cover -html=coverage.out
```

## Roadmap

- [x] Create cache structure
- [x] Implement `Set`
- [x] Implement `Get`
- [x] Implement `Delete`
- [x] Implement `Has`
- [x] Implement `Len`
- [x] Add lazy TTL expiration on read
- [x] Add background cleanup
- [x] Add tests
- [x] Add public API documentation
- [x] Add usage examples
- [ ] Add CI quality checks
- [ ] Add benchmarks
- [ ] Add statistics
- [ ] Add hooks/callbacks
- [ ] Explore generics support

## License

MIT