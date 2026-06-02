# tinyttl

A lightweight in-memory TTL cache for Go.

## Features

- In-memory key-value storage
- TTL expiration per entry
- Thread-safe access
- Lazy expiration on read
- Optional background cleanup
- Simple and minimal API

## Status

`tinyttl` is currently in early development. The first goal is to provide a small, reliable, and easy-to-use TTL cache for Go applications.

## Planned API

```go
cache := tinyttl.New()

_ = cache.Set("user:1", "Diego", 5*time.Minute)

value, ok := cache.Get("user:1")
if ok {
    fmt.Println(value)
}
```

## Roadmap

### MVP
- [ ] Create cache structure
- [ ] Implement `Set`
- [ ] Implement `Get`
- [ ] Implement `Delete`
- [ ] Implement `Has`
- [ ] Implement `Len`
- [ ] Add TTL expiration
- [ ] Add background cleanup
- [ ] Add tests

### Next
- [ ] Add benchmarks
- [ ] Add generics support
- [ ] Add stats
- [ ] Add eviction callbacks

## Motivation

This project is being built to deepen Go knowledge through practical work on:
- concurrency
- synchronization
- API design
- testing
- benchmarking

## License

MIT
