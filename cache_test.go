package tinyttl

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	cache := New()
	if cache == nil {
		t.Fatal("expected New() to return a non-nil Cache instance")
	}
	if cache.items == nil {
		t.Fatal("expected Cache.items to be initialized")
	}
}

func TestCache_SetAndGet(t *testing.T) {
	cache := New()

	cache.Set("key1", "value1", time.Minute)

	value, found := cache.Get("key1")
	if !found {
		t.Fatal("expected to find key1 in cache")
	}
	if value != "value1" {
		t.Fatalf("expected value1, got %v", value)
	}
}

func TestCache_ReturnsFalseForMissingKey(t *testing.T) {
	cache := New()

	value, found := cache.Get("missingKey")
	if found {
		t.Fatal("expected not to find missingKey in cache")
	}
	if value != nil {
		t.Fatalf("expected nil value for missingKey, got %v", value)
	}
}

func TestCache_OverwritesExistingKey(t *testing.T) {
	cache := New()

	cache.Set("key1", "value1", time.Minute)
	cache.Set("key1", "value2", time.Minute)

	value, found := cache.Get("key1")
	if !found {
		t.Fatal("expected to find key1 in cache")
	}
	if value != "value2" {
		t.Fatalf("expected value2, got %v", value)
	}
}

func TestCache_GetReturnsFalseForExpiredKey(t *testing.T) {
	cache := New()

	cache.Set("key1", "value1", time.Millisecond*100)
	time.Sleep(time.Millisecond * 200)

	value, found := cache.Get("key1")
	if found {
		t.Fatal("expected not to find key1 in cache after expiration")
	}
	if value != nil {
		t.Fatalf("expected nil value for expired key1, got %v", value)
	}
}

func TestCache_ExpiredKeyIsDeletedOnGet(t *testing.T) {
	cache := New()

	cache.Set("key1", "value1", 10*time.Millisecond)
	time.Sleep(20 * time.Millisecond)

	_, _ = cache.Get("key1")

	if _, exists := cache.items["key1"]; exists {
		t.Fatal("expected expired key to be deleted from cache")
	}
}
