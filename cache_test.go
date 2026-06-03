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

func TestCache_DeleteRemovesExistingKey(t *testing.T) {
	cache := New()

	cache.Set("ABC", "random value", time.Minute)
	cache.Delete("ABC")

	value, found := cache.Get("ABC")
	if found {
		t.Fatal("expected not to find ABC in cache after deletion")
	}
	if value != nil {
		t.Fatalf("expected nil value for deleted key ABC, got %v", value)
	}
}

func TestCache_DeleteMissingKeyDoesNothing(t *testing.T) {
	cache := New()

	cache.Delete("nonExistentKey")

	if cache.Len() != 0 {
		t.Fatalf("expected cache to be empty after deleting non-existent key, got length %d", cache.Len())
	}
}

func TestCache_HasReturnsTrueForExistingKey(t *testing.T) {
	cache := New()

	cache.Set("existingKey", "value", time.Minute)

	if !cache.Has("existingKey") {
		t.Fatal("expected Has to return true for existingKey")
	}
}

func TestCache_HasReturnsFalseForMissingKey(t *testing.T) {
	cache := New()

	if cache.Has("nonExistentKey") {
		t.Fatal("expected Has to return false for nonExistentKey")
	}
}

func TestCache_HasReturnsFalseForExpiredKey(t *testing.T) {
	cache := New()

	cache.Set("tempKey", "value", time.Millisecond*10)
	time.Sleep(20 * time.Millisecond)

	if cache.Has("tempKey") {
		t.Fatal("expected Has to return false for expired tempKey")
	}
}

func TestCache_HasDeletesExpiredKey(t *testing.T) {
	cache := New()

	cache.Set("tempKey", "value", time.Millisecond*10)
	time.Sleep(20 * time.Millisecond)

	_ = cache.Has("tempKey")

	if _, exists := cache.items["tempKey"]; exists {
		t.Fatal("expected expired key to be deleted from cache after Has check")
	}
}

func TestCache_LenReturnsNumberOfNonExpiredItems(t *testing.T) {
	cache := New()

	cache.Set("key1", "value1", time.Minute)
	cache.Set("key3", "value3", time.Minute)

	if cache.Len() != 2 {
		t.Fatalf("expected Len to return 2, got %d", cache.Len())
	}
}

func TestCache_LenDoesNotCountExpiredItems(t *testing.T) {
	cache := New()

	cache.Set("key1", "value1", time.Millisecond*10)
	cache.Set("key2", "value2", time.Minute)
	cache.Set("key3", "value3", time.Minute)
	cache.Set("key4", "value3", time.Hour)

	time.Sleep(20 * time.Millisecond)

	if cache.Len() != 3 {
		t.Fatalf("expected cache length to be 3 after expiration, got %d", cache.Len())
	}
}

func TestCache_LenDeletesExpiredItems(t *testing.T) {
	cache := New()

	cache.Set("toExpire", "value1", 10*time.Millisecond)
	cache.Set("validItem", "value2", time.Minute)

	time.Sleep(20 * time.Millisecond)

	count := cache.Len()

	if count != 1 {
		t.Fatalf("expected cache length to be 1, got %d", count)
	}

	if cache.Has("toExpire") {
		t.Fatalf("expected expired key to be removed from cache")
	}

	if !cache.Has("validItem") {
		t.Fatalf("expected non-expired item to remain in cache")
	}
}
