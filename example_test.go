package tinyttl_test

import (
	"fmt"
	"time"

	"github.com/diegofer11/tinyttl"
)

func ExampleNew() {
	cache := tinyttl.New()

	cache.Set("message", "Hello", time.Minute)

	value, found := cache.Get("message")
	if !found {
		fmt.Println("not found")
		return
	}

	fmt.Println(value)
	// Output: Hello
}

func ExampleCache_Set_noExpiration() {
	cache := tinyttl.New()

	cache.Set("config", "enabled", 0)

	value, found := cache.Get("config")
	if !found {
		fmt.Println("not found")
		return
	}

	fmt.Println(value)
	// Output: enabled
}

func ExampleWithCleanupInterval() {
	cache := tinyttl.New(tinyttl.WithCleanupInterval(10 * time.Millisecond))
	defer cache.Close()

	cache.Set("session", "random session 123", time.Minute)

	if cache.Has("session") {
		fmt.Println("session exists")
	}

	// Output: session exists
}

func ExampleCache_Stats() {
	cache := tinyttl.New()

	cache.Set("message", "hello", time.Minute)
	_, _ = cache.Get("message")
	_, _ = cache.Get("missing")

	stats := cache.Stats()

	fmt.Printf("hits=%d misses=%d sets=%d\n", stats.Hits, stats.Misses, stats.Sets)
	// Output: hits=1 misses=1 sets=1
}

func ExampleWithHooks() {
	cache := tinyttl.New(
		tinyttl.WithHooks(tinyttl.Hooks{
			OnSet: func(key string, value any) {
				fmt.Printf("set %s=%v\n", key, value)
			},
		}),
	)

	cache.Set("user:1", "Diego", time.Minute)

	// Output: set user:1=Diego
}
