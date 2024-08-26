package benchmark

import (
	"testing"

	"github.com/lovelysunlight/lru-go"
	"github.com/pioz/faker"
)

func BenchmarkCache_Rand(b *testing.B) {
	l, err := lru.New[int64, int64](8192, lru.DisableDeepCopy[int64, int64]())
	if err != nil {
		b.Fatalf("err: %v", err)
	}

	trace := make([]int64, b.N*2)
	for i := 0; i < b.N*2; i++ {
		trace[i] = faker.Int64() % 32768
	}
	b.ResetTimer()
	var hit, miss int
	for i := 0; i < 2*b.N; i++ {
		if i%2 == 0 {
			l.Put(trace[i], trace[i])
		} else {
			if _, ok := l.Get(trace[i]); ok {
				hit++
			} else {
				miss++
			}
		}
	}
	b.Logf("hit: %d miss: %d ratio: %f", hit, miss, float64(hit)/float64(hit+miss))
}

func BenchmarkCache_Freq(b *testing.B) {
	l, err := lru.New[int64, int64](8192, lru.DisableDeepCopy[int64, int64]())
	if err != nil {
		b.Fatalf("err: %v", err)
	}

	trace := make([]int64, b.N*2)
	for i := 0; i < b.N*2; i++ {
		if i%2 == 0 {
			trace[i] = faker.Int64() % 16384
		} else {
			trace[i] = faker.Int64() % 32768
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		l.Put(trace[i], trace[i])
	}
	var hit, miss int
	for i := 0; i < b.N; i++ {
		if _, ok := l.Get(trace[i]); ok {
			hit++
		} else {
			miss++
		}
	}
	b.Logf("hit: %d miss: %d ratio: %f", hit, miss, float64(hit)/float64(hit+miss))
}
