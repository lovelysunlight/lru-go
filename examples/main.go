package main

import (
	"fmt"

	"github.com/lovelysunlight/lru-go"
)

func main() {
	cache, _ := lru.New[string, string](2)
	cache.Put("apple", "red")
	cache.Put("banana", "yellow")

	var (
		r, v string
		ok   bool
	)

	r, ok = cache.Get("apple")
	fmt.Printf("Get() found: %v, value: %q\n", ok, r)

	r, ok = cache.Get("banana")
	fmt.Printf("Get() found: %v, value: %q\n", ok, r)

	r, ok = cache.Get("pear")
	fmt.Printf("Get() found: %v, value: %q\n", ok, r)

	r, ok = cache.Peek("apple")
	fmt.Printf("Peek() found: %v, value: %q\n", ok, r)

	r, ok = cache.Peek("banana")
	fmt.Printf("Peek() found: %v, value: %q\n", ok, r)

	r, ok = cache.Peek("pear")
	fmt.Printf("Peek() found: %v, value: %q\n", ok, r)

	r, ok = cache.Remove("banana")
	fmt.Printf("Remove() found: %v, value: %q\n", ok, r)

	r, v, ok = cache.RemoveOldest()
	fmt.Printf("RemoveOldest() found: %v, key: %q, value: %q\n", ok, r, v)

	fmt.Printf("Len() = : %v\n", cache.Len())
	fmt.Printf("Cap() = : %v\n", cache.Cap())
}
