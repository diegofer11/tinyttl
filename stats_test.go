package tinyttl

import (
	"testing"
	"time"
)

func TestCache_Stats_InitialValues(t *testing.T) {
	cache := New()

	stats := cache.Stats()

	if stats.Sets != 0 {
		t.Fatalf("expected Sets to be 0, got %d", stats.Sets)
	}
	if stats.Hits != 0 {
		t.Fatalf("expected Hits to be 0, got %d", stats.Hits)
	}
	if stats.Misses != 0 {
		t.Fatalf("expected Misses to be 0, got %d", stats.Misses)
	}
	if stats.Deletes != 0 {
		t.Fatalf("expected Deletes to be 0, got %d", stats.Deletes)
	}
	if stats.Expirations != 0 {
		t.Fatalf("expected Expirations to be 0, got %d", stats.Expirations)
	}
}

func TestCache_Stats_Set(t *testing.T) {
	cache := New()

	cache.Set("key1", "value1", time.Minute)
	cache.Set("key2", "value2", time.Minute)

	stats := cache.Stats()

	if stats.Sets != 2 {
		t.Fatalf("expected Sets to be 2, got %d", stats.Sets)
	}
}

func TestCache_Stats_GetHit(t *testing.T) {
	cache := New()
	cache.Set("key", "value", time.Minute)

	value, found := cache.Get("key")
	if !found {
		t.Fatal("expected key to be found")
	}
	if value != "value" {
		t.Fatalf("expected value to be %q, got %v", "value", value)
	}

	stats := cache.Stats()

	if stats.Hits != 1 {
		t.Fatalf("expected Hits to be 1, got %d", stats.Hits)
	}
	if stats.Misses != 0 {
		t.Fatalf("expected Misses to be 0, got %d", stats.Misses)
	}
}

func TestCache_Stats_GetMiss(t *testing.T) {
	cache := New()

	value, found := cache.Get("missing")
	if found {
		t.Fatal("expected key to be missing")
	}
	if value != nil {
		t.Fatalf("expected nil value, got %v", value)
	}

	stats := cache.Stats()

	if stats.Misses != 1 {
		t.Fatalf("expected Misses to be 1, got %d", stats.Misses)
	}
	if stats.Hits != 0 {
		t.Fatalf("expected Hits to be 0, got %d", stats.Hits)
	}
}

func TestCache_Stats_DeleteExisting(t *testing.T) {
	cache := New()
	cache.Set("key", "value", time.Minute)

	cache.Delete("key")

	stats := cache.Stats()

	if stats.Deletes != 1 {
		t.Fatalf("expected Deletes to be 1, got %d", stats.Deletes)
	}
}

func TestCache_Stats_DeleteMissing(t *testing.T) {
	cache := New()

	cache.Delete("missing")

	stats := cache.Stats()

	if stats.Deletes != 0 {
		t.Fatalf("expected Deletes to be 0, got %d", stats.Deletes)
	}
}

func TestCache_Stats_ExpirationOnGet(t *testing.T) {
	cache := New()
	cache.Set("key", "value", 10*time.Millisecond)

	time.Sleep(20 * time.Millisecond)

	value, found := cache.Get("key")
	if found {
		t.Fatal("expected key to be expired")
	}
	if value != nil {
		t.Fatalf("expected nil value, got %v", value)
	}

	stats := cache.Stats()

	if stats.Expirations != 1 {
		t.Fatalf("expected Expirations to be 1, got %d", stats.Expirations)
	}
	if stats.Misses != 1 {
		t.Fatalf("expected Misses to be 1, got %d", stats.Misses)
	}
}

func TestCache_Stats_HasDoesNotCountHitOrMiss(t *testing.T) {
	cache := New()
	cache.Set("key", "value", time.Minute)

	if !cache.Has("key") {
		t.Fatal("expected key to exist")
	}

	stats := cache.Stats()

	if stats.Hits != 0 {
		t.Fatalf("expected Hits to be 0, got %d", stats.Hits)
	}
	if stats.Misses != 0 {
		t.Fatalf("expected Misses to be 0, got %d", stats.Misses)
	}
}

func TestCache_Stats_ExpirationOnHas(t *testing.T) {
	cache := New()
	cache.Set("key", "value", 10*time.Millisecond)

	time.Sleep(20 * time.Millisecond)

	if cache.Has("key") {
		t.Fatal("expected key to be expired")
	}

	stats := cache.Stats()

	if stats.Expirations != 1 {
		t.Fatalf("expected Expirations to be 1, got %d", stats.Expirations)
	}
}

func TestCache_Stats_ExpirationOnLen(t *testing.T) {
	cache := New()
	cache.Set("expired", "value1", 10*time.Millisecond)
	cache.Set("valid", "value2", time.Minute)

	time.Sleep(20 * time.Millisecond)

	count := cache.Len()
	if count != 1 {
		t.Fatalf("expected Len to return 1, got %d", count)
	}

	stats := cache.Stats()

	if stats.Expirations != 1 {
		t.Fatalf("expected Expirations to be 1, got %d", stats.Expirations)
	}
}

func TestCache_Stats_ExpirationOnBackgroundCleanup(t *testing.T) {
	cache := New(WithCleanupInterval(10 * time.Millisecond))
	defer cache.Close()

	cache.Set("key", "value", 10*time.Millisecond)

	time.Sleep(40 * time.Millisecond)

	stats := cache.Stats()

	if stats.Expirations != 1 {
		t.Fatalf("expected Expirations to be 1, got %d", stats.Expirations)
	}
}
