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

func BenchmarkCache_LenBySize(b *testing.B) {
	sizes := []int{100, 1000, 10000}

	for _, size := range sizes {
		b.Run(strconv.Itoa(size), func(b *testing.B) {
			cache := New()

			for i := range size {
				cache.Set(strconv.Itoa(i), "value", time.Minute)
			}

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_ = cache.Len()
			}
		})
	}
}

func BenchmarkCache_GetBySize(b *testing.B) {
	sizes := []int{100, 1000, 10000}

	for _, size := range sizes {
		b.Run(strconv.Itoa(size), func(b *testing.B) {
			cache := New()

			for i := range size {
				cache.Set(strconv.Itoa(i), "value", time.Minute)
			}

			targetKey := strconv.Itoa(size / 2)

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_, _ = cache.Get(targetKey)
			}
		})
	}
}

func BenchmarkCache_MixedReadWrite(b *testing.B) {
	cache := New()

	for i := range 1000 {
		cache.Set(strconv.Itoa(i), "value", time.Minute)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		switch i % 20 {
		case 0, 1:
			cache.Set(strconv.Itoa(i+1000), "value", time.Minute)
		case 2:
			cache.Delete(strconv.Itoa(i % 1000))
		case 3, 4:
			_ = cache.Has(strconv.Itoa(i % 1000))
		default:
			_, _ = cache.Get(strconv.Itoa(i % 1000))
		}
	}
}

func BenchmarkCache_GetWithExpiredEntries(b *testing.B) {
	cache := New()

	for i := range 1000 {
		cache.Set("valid-"+strconv.Itoa(i), "value", time.Minute)
	}

	for i := range 100 {
		cache.Set("expired-"+strconv.Itoa(i), "value", 10*time.Millisecond)
	}

	time.Sleep(20 * time.Millisecond)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if i%10 == 0 {
			_, _ = cache.Get("expired-" + strconv.Itoa(i%100))
		} else {
			_, _ = cache.Get("valid-" + strconv.Itoa(i%1000))
		}
	}
}

func BenchmarkCache_LenWithExpiredEntries(b *testing.B) {
	for _, size := range []int{100, 1000, 10000} {
		b.Run(strconv.Itoa(size), func(b *testing.B) {
			cache := New()

			expiredCount := size / 10
			validCount := size - expiredCount

			for i := range validCount {
				cache.Set("valid-"+strconv.Itoa(i), "value", time.Minute)
			}

			for i := range expiredCount {
				cache.Set("expired-"+strconv.Itoa(i), "value", 10*time.Millisecond)
			}

			time.Sleep(20 * time.Millisecond)

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_ = cache.Len()
			}
		})
	}
}

func BenchmarkCache_GetParallel(b *testing.B) {
	cache := New()
	cache.Set("key", "value", time.Minute)

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = cache.Get("key")
		}
	})
}
