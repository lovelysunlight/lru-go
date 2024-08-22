# LRU Cache

[![Build Badge]][build status]
[![Codecov Badge]][coverage status]
[![License Badge]][license]

An implementation of a LRU cache. The cache supports `Push`, `Put`, `Get` `Peek` and `Pop` operations,
all of which are O(1). This package was heavily influenced by the [LRU Cache implementation in a Rust crate].

## Example

Below is a simple example of how to instantiate and use a LRU cache.

```golang
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

	r, ok = cache.Pop("banana")
	fmt.Printf("Pop() found: %v, value: %q\n", ok, r)

	r, v, ok = cache.RemoveOldest()
	fmt.Printf("RemoveOldest() found: %v, key: %q, value: %q\n", ok, r, v)

	fmt.Printf("Len() = : %v\n", cache.Len())
	fmt.Printf("Cap() = : %v\n", cache.Cap())
}
```

[build badge]: https://github.com/lovelysunlight/lru-go/actions/workflows/ci.yaml/badge.svg
[build status]: https://github.com/lovelysunlight/lru-go/actions/workflows/ci.yaml
[codecov badge]: https://codecov.io/gh/lovelysunlight/lru-go/branch/master/graph/badge.svg
[coverage status]: https://codecov.io/gh/lovelysunlight/lru-go
[license badge]: https://img.shields.io/badge/license-MIT-blue.svg
[license]: https://raw.githubusercontent.com/lovelysunlight/lru-go/master/LICENSE
[LRU Cache implementation in a Rust crate]: https://github.com/jeromefroe/lru-rs
