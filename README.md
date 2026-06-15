# tinyttl

[![Go Reference](https://pkg.go.dev/badge/github.com/diegofer11/tinyttl.svg)](https://pkg.go.dev/github.com/diegofer11/tinyttl)
[![Go Report Card](https://goreportcard.com/badge/github.com/diegofer11/tinyttl)](https://goreportcard.com/report/github.com/diegofer11/tinyttl)

`tinyttl` is a lightweight in-memory TTL cache for Go.

It provides a small and stable API for storing values with optional expiration times, lazy expiration on access, optional background cleanup, statistics, and event hooks.

## Stability

`tinyttl` `v1.0.0` marks the first stable release of the library.

The core public API is considered stable and suitable for adoption in production applications that need a small in-memory TTL cache with predictable behavior and minimal setup.

## Features

- In-memory cache with string keys
- Store arbitrary values
- Per-item TTL support
- Lazy expiration on `Get`, `Has`, and `Len`
- Optional background cleanup
- Cache statistics
- Event hooks
- Small and focused API

## When to use

`tinyttl` is a good fit when you need:

- a small in-memory cache inside a Go process
- TTL-based expiration
- a minimal API with low setup overhead
- lightweight observability through stats and hooks
- a simple cache for services, workers, CLI tools, or internal libraries

## What it does not do

`tinyttl` is intentionally small and focused. It does not currently provide:

- distributed caching
- persistence
- capacity-based eviction
- LRU or LFU policies
- sharding
- advanced expiration policies such as expire-after-access
- generics-based typed APIs

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

A TTL greater than zero sets an expiration time for the item.

A TTL less than or equal to zero means the item does not expire.

Expired items are removed lazily when accessed through:

- `Get`
- `Has`
- `Len`

## Semantics

- A TTL greater than zero sets an expiration time for the item.
- A TTL less than or equal to zero means the item does not expire.
- Expiration is lazy on `Get`, `Has`, and `Len`.
- Optional background cleanup can remove expired items periodically.
- `Has` does not count as a cache hit or miss.
- `Delete` only affects statistics when an existing key is removed.
- If background cleanup is enabled, call `Close()` when the cache is no longer needed.

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

### Read cache statistics

```go
stats := cache.Stats()
fmt.Printf("hits=%d misses=%d sets=%d\n", stats.Hits, stats.Misses, stats.Sets)
```

### Register hooks

```go
cache := tinyttl.New(
	tinyttl.WithHooks(tinyttl.Hooks{
		OnSet: func(key string, value any) {
			fmt.Printf("set %s\n", key)
		},
		OnDelete: func(key string, value any) {
			fmt.Printf("delete %s\n", key)
		},
		OnExpire: func(key string, value any) {
			fmt.Printf("expired %s\n", key)
		},
	}),
)
```

## Development

Run tests:

```bash
go test
```

Run tests with race detection:

```bash
go test -race
```

Run tests with coverage:

```bash
go test -coverprofile=coverage.out
go tool cover -func=coverage.out
```

Open the HTML coverage report:

```bash
go tool cover -html=coverage.out
```

Run static analysis:

```bash
go vet
```

## Benchmarks

Run benchmarks with:

```bash
go test -run=^$ -bench=. -benchmem
```

## Roadmap

Future releases may expand functionality, but `v1.0.0` establishes the core API and behavior as stable.

Possible future areas of exploration include:

- generics support
- sharding
- advanced expiration policies
- additional performance optimizations

## License

MIT