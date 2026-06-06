package tinyttl

import (
	"strconv"
	"testing"
	"time"
)

func BenchmarkCache_GetParallel_SharedKey(b *testing.B) {
	cache := New()
	cache.Set("shared", "value", time.Minute)

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = cache.Get("shared")
		}
	})
}

func BenchmarkCache_GetParallel_DistributedKeys(b *testing.B) {
	cache := New()

	keys := make([]string, 1024)
	for i := range keys {
		key := strconv.Itoa(i)
		keys[i] = key
		cache.Set(key, "value", time.Minute)
	}

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			_, _ = cache.Get(keys[i%len(keys)])
			i++
		}
	})
}

func BenchmarkCache_SetParallel(b *testing.B) {
	cache := New()

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			cache.Set(strconv.Itoa(i), "value", time.Minute)
			i++
		}
	})
}

func BenchmarkCache_DeleteParallel(b *testing.B) {
	cache := New()

	keys := make([]string, 1024)
	for i := range keys {
		key := strconv.Itoa(i)
		keys[i] = key
		cache.Set(key, "value", time.Minute)
	}

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			cache.Delete(keys[i%len(keys)])
			i++
		}
	})
}

func BenchmarkCache_MixedParallel_ReadHeavy(b *testing.B) {
	cache := New()

	keys := make([]string, 1024)
	for i := range keys {
		key := strconv.Itoa(i)
		keys[i] = key
		cache.Set(key, "value", time.Minute)
	}

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := keys[i%len(keys)]

			switch i % 20 {
			case 0:
				cache.Set(key, "value", time.Minute)
			case 1:
				cache.Delete(key)
			case 2, 3:
				_ = cache.Has(key)
			default:
				_, _ = cache.Get(key)
			}
			i++
		}
	})
}

func BenchmarkCache_MixedParallel_WriteHeavy(b *testing.B) {
	cache := New()

	keys := make([]string, 1024)
	for i := range keys {
		key := strconv.Itoa(i)
		keys[i] = key
		cache.Set(key, "value", time.Minute)
	}

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := keys[i%len(keys)]

			switch i % 10 {
			case 0, 1, 2, 3:
				cache.Set(key, "value", time.Minute)
			case 4, 5:
				cache.Delete(key)
			case 6:
				_ = cache.Has(key)
			default:
				_, _ = cache.Get(key)
			}
			i++
		}
	})
}
