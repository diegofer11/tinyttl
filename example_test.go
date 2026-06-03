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
