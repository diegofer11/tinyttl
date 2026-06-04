package tinyttl

import (
	"strconv"
	"testing"
	"time"
)

func BenchmarkCache_Set(b *testing.B) {
	cache := New()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cache.Set(strconv.Itoa(i), "value", time.Minute)
	}
}

func BenchmarkCache_Get(b *testing.B) {
	cache := New()
	cache.Set("key", "value", time.Minute)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = cache.Get("key")
	}
}

func BenchmarkCache_GetMissing(b *testing.B) {
	cache := New()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = cache.Get("missing")
	}
}

func BenchmarkCache_Has(b *testing.B) {
	cache := New()
	cache.Set("key", "value", time.Minute)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = cache.Has("key")
	}
}

func BenchmarkCache_Delete(b *testing.B) {
	cache := New()

	for i := 0; i < b.N; i++ {
		cache.Set(strconv.Itoa(i), "value", time.Minute)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cache.Delete(strconv.Itoa(i))
	}
}

func BenchmarkCache_Len(b *testing.B) {
	cache := New()

	for i := range 100 {
		cache.Set(strconv.Itoa(i), "value", time.Minute)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = cache.Len()
	}
}
