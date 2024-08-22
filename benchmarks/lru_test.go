package benchmark

import (
	"testing"

	"github.com/lovelysunlight/lru-go"
	"github.com/pioz/faker"
)

type Object struct {
	Str1  string
	Str2  string
	Str3  string
	Str4  string
	Age1  int
	Age2  int
	Age3  int
	Age4  int
	Bool1 bool
	Bool2 bool
	Bool3 bool
	Bool4 bool

	Obj1 *Object
	Obj2 *Object
	Obj3 *Object
	Obj4 *Object

	Map1 map[string]string
	Map2 map[string]int
	Map3 map[string]bool
	Map4 map[string]float32
	Map5 map[string]*Object
}

func BenchmarkLruCache_Put(b *testing.B) {
	b.Run("int", func(b *testing.B) {
		cache, _ := lru.New[int, int](10)
		for i := 0; i < b.N; i++ {
			cache.Put(i, i)
		}
	})

	b.Run("float", func(b *testing.B) {
		cache, _ := lru.New[int, float32](10)
		for i := 0; i < b.N; i++ {
			cache.Push(i, 1.1)
		}
	})

	b.Run("bool", func(b *testing.B) {
		cache, _ := lru.New[int, bool](10)
		for i := 0; i < b.N; i++ {
			cache.Push(i, true)
		}
	})

	b.Run("string", func(b *testing.B) {
		cache, _ := lru.New[int, string](10)
		for i := 0; i < b.N; i++ {
			cache.Push(i, "aaaaa")
		}
	})

	b.Run("struct", func(b *testing.B) {
		cache, _ := lru.New[int, *Object](10)
		for i := 0; i < b.N; i++ {
			cache.Push(i, &Object{
				Str1:  faker.String(),
				Str2:  faker.String(),
				Str3:  faker.String(),
				Str4:  faker.String(),
				Age1:  faker.Int(),
				Age2:  faker.Int(),
				Age3:  faker.Int(),
				Age4:  faker.Int(),
				Bool1: faker.Bool(),
				Bool2: faker.Bool(),
				Bool3: faker.Bool(),
				Bool4: faker.Bool(),

				Obj1: &Object{},
				Obj2: &Object{},
				Obj3: &Object{},
				Obj4: nil,

				Map1: map[string]string{
					faker.String(): faker.String(),
					faker.String(): faker.String(),
					faker.String(): faker.String(),
					faker.String(): faker.String(),
				},
				Map2: map[string]int{
					faker.String(): faker.Int(),
				},
				Map3: map[string]bool{
					faker.String(): faker.Bool(),
				},
				Map4: map[string]float32{
					faker.String(): faker.Float32(),
				},
			})
		}
	})
}

func BenchmarkLruCache_Get(b *testing.B) {
	b.Run("int", func(b *testing.B) {
		cache, _ := lru.New[int, int](10)
		cache.Put(1, 1)
		for i := 0; i < b.N; i++ {
			cache.Get(1)
		}
	})

	b.Run("float", func(b *testing.B) {
		cache, _ := lru.New[int, float32](10)
		cache.Put(1, 1.1)
		for i := 0; i < b.N; i++ {
			cache.Get(1)
		}
	})

	b.Run("bool", func(b *testing.B) {
		cache, _ := lru.New[int, bool](10)
		cache.Put(1, true)
		for i := 0; i < b.N; i++ {
			cache.Get(1)
		}
	})

	b.Run("string", func(b *testing.B) {
		cache, _ := lru.New[int, string](10)
		cache.Put(1, "aaaaa")
		for i := 0; i < b.N; i++ {
			cache.Get(1)
		}
	})

	b.Run("[]int", func(b *testing.B) {
		b.Run("immutable", func(b *testing.B) {
			cache, _ := lru.New(10, lru.WithImmutable[int, []int]())
			cache.Put(1, []int{1, 2, 3, 4, 5, 6})
			for i := 0; i < b.N; i++ {
				cache.Get(1)
			}
		})
		b.Run("mutable", func(b *testing.B) {
			cache, _ := lru.New(10, lru.WithMutable[int, []int]())
			cache.Put(1, []int{1, 2, 3, 4, 5, 6})
			for i := 0; i < b.N; i++ {
				cache.Get(1)
			}
		})
	})
}

func BenchmarkLruCache_Peek(b *testing.B) {
	b.Run("int", func(b *testing.B) {
		cache, _ := lru.New[int, int](10)
		cache.Put(1, 1)
		for i := 0; i < b.N; i++ {
			cache.Peek(1)
		}
	})

	b.Run("float", func(b *testing.B) {
		cache, _ := lru.New[int, float32](10)
		cache.Put(1, 1.1)
		for i := 0; i < b.N; i++ {
			cache.Peek(1)
		}
	})

	b.Run("bool", func(b *testing.B) {
		cache, _ := lru.New[int, bool](10)
		cache.Put(1, true)
		for i := 0; i < b.N; i++ {
			cache.Peek(1)
		}
	})

	b.Run("string", func(b *testing.B) {
		cache, _ := lru.New[int, string](10)
		cache.Put(1, "aaaaa")
		for i := 0; i < b.N; i++ {
			cache.Peek(1)
		}
	})
}

func BenchmarkLruCache_Pop(b *testing.B) {
	len := 10000000
	b.Run("int", func(b *testing.B) {
		cache, _ := lru.New[int, int](len)
		for i := 0; i < len; i++ {
			cache.Put(i, i)
		}
		b.Run("benchmark", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				cache.Pop(1)
			}
		})
	})

	b.Run("float", func(b *testing.B) {
		cache, _ := lru.New[int, float32](len)
		for i := 0; i < len; i++ {
			cache.Put(i, 1.1)
		}
		b.Run("benchmark", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				cache.Pop(1)
			}
		})
	})

	b.Run("bool", func(b *testing.B) {
		cache, _ := lru.New[int, bool](len)
		for i := 0; i < len; i++ {
			cache.Put(i, true)
		}
		b.Run("benchmark", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				cache.Pop(1)
			}
		})
	})

	b.Run("string", func(b *testing.B) {
		cache, _ := lru.New[int, string](len)
		for i := 0; i < len; i++ {
			cache.Put(i, "aaaaa")
		}
		b.Run("benchmark", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				cache.Pop(1)
			}
		})
	})
}
