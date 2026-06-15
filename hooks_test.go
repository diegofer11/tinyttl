package tinyttl

import (
	"sync"
	"testing"
	"time"
)

func TestCache_OnSetHook(t *testing.T) {
	var (
		called bool
		gotKey string
		gotVal any
	)

	cache := New(WithHooks(Hooks{
		OnSet: func(key string, value any) {
			called = true
			gotKey = key
			gotVal = value
		},
	}))

	cache.Set("user:1", "Diego", time.Minute)

	if !called {
		t.Fatal("expected OnSet hook to be called")
	}
	if gotKey != "user:1" {
		t.Fatalf("expected key to be user:1, got %s", gotKey)
	}
	if gotVal != "Diego" {
		t.Fatalf("expected value to be Diego, got %v", gotVal)
	}
}

func TestCache_OnDeleteHook(t *testing.T) {
	var (
		called bool
		gotKey string
		gotVal any
	)

	cache := New(WithHooks(Hooks{
		OnDelete: func(key string, value any) {
			called = true
			gotKey = key
			gotVal = value
		},
	}))

	cache.Set("user:1", "Diego", time.Minute)
	cache.Delete("user:1")

	if !called {
		t.Fatal("expected OnDelete hook to be called")
	}
	if gotKey != "user:1" {
		t.Fatalf("expected key to be user:1, got %s", gotKey)
	}
	if gotVal != "Diego" {
		t.Fatalf("expected value to be Diego, got %v", gotVal)
	}
}

func TestCache_OnDeleteHookNotCalledForMissingKey(t *testing.T) {
	called := false

	cache := New(WithHooks(Hooks{
		OnDelete: func(key string, value any) {
			called = true
		},
	}))

	cache.Delete("missing")

	if called {
		t.Fatal("expected OnDelete hook not to be called")
	}
}

func TestCache_OnExpireHookOnGet(t *testing.T) {
	var (
		called bool
		gotKey string
		gotVal any
	)

	cache := New(WithHooks(Hooks{
		OnExpire: func(key string, value any) {
			called = true
			gotKey = key
			gotVal = value
		},
	}))

	cache.Set("session", "abc123", 10*time.Millisecond)
	time.Sleep(20 * time.Millisecond)

	_, _ = cache.Get("session")

	if !called {
		t.Fatal("expected OnExpire hook to be called")
	}
	if gotKey != "session" {
		t.Fatalf("expected key to be session, got %s", gotKey)
	}
	if gotVal != "abc123" {
		t.Fatalf("expected value to be abc123, got %v", gotVal)
	}
}

func TestCache_OnExpireHookOnHas(t *testing.T) {
	called := false

	cache := New(WithHooks(Hooks{
		OnExpire: func(key string, value any) {
			called = true
		},
	}))

	cache.Set("session", "abc123", 10*time.Millisecond)
	time.Sleep(20 * time.Millisecond)

	_ = cache.Has("session")

	if !called {
		t.Fatal("expected OnExpire hook to be called")
	}
}

func TestCache_OnExpireHookOnLen(t *testing.T) {
	var called int

	cache := New(WithHooks(Hooks{
		OnExpire: func(key string, value any) {
			called++
		},
	}))

	cache.Set("expired1", "a", 10*time.Millisecond)
	cache.Set("expired2", "b", 10*time.Millisecond)
	cache.Set("valid", "c", time.Minute)

	time.Sleep(20 * time.Millisecond)

	if cache.Len() != 1 {
		t.Fatalf("expected Len to return 1")
	}

	if called != 2 {
		t.Fatalf("expected OnExpire hook to be called 2 times, got %d", called)
	}
}

func TestCache_OnExpireHookOnBackgroundCleanup(t *testing.T) {
	var (
		mu     sync.Mutex
		called int
	)

	cache := New(
		WithCleanupInterval(10*time.Millisecond),
		WithHooks(Hooks{
			OnExpire: func(key string, value any) {
				mu.Lock()
				called++
				mu.Unlock()
			},
		}),
	)
	defer cache.Close()

	cache.Set("session", "abc123", 10*time.Millisecond)

	time.Sleep(40 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	if called != 1 {
		t.Fatalf("expected OnExpire hook to be called 1 time, got %d", called)
	}
}

func TestCache_OnSetHookOnOverwrite(t *testing.T) {
	var called int

	cache := New(WithHooks(Hooks{
		OnSet: func(key string, value any) {
			called++
		},
	}))

	cache.Set("key", "value", time.Minute)
	cache.Set("key", "new-value", time.Minute)

	if called != 2 {
		t.Fatalf("expected OnSet hook to be called 2 times, got %d", called)
	}
}

func TestCache_OnExpireHookIsCalledOnlyOncePerEntry(t *testing.T) {
	var called int

	cache := New(WithHooks(Hooks{
		OnExpire: func(key string, value any) {
			called++
		},
	}))

	cache.Set("key", "value", 10*time.Millisecond)
	time.Sleep(20 * time.Millisecond)

	_, _ = cache.Get("key")
	_ = cache.Has("key")
	_ = cache.Len()

	if called != 1 {
		t.Fatalf("expected OnExpire hook to be called 1 time, got %d", called)
	}
}
