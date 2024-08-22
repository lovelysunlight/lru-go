# LRU Cache

[![Build Badge]][build status]
[![Codecov Badge]][coverage status]
[![License Badge]][license]

An implementation of a LRU cache. The cache supports `Put`, `Get` `Peek` and `Pop` operations,
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
    cache := lru.New[string, string](2)
    cache.Put("apple", "red")
    cache.Put("banana", "yellow")

    fmt.Println(cache.get("apple").Unwrap()) // "red"
    fmt.Println(cache.get("banana").Unwrap()) // "yellow"
    fmt.Println(cache.get("pear").IsSome()) // false

    fmt.Println(cache.put("banana", "foo").Unwrap()) // "yellow"
    fmt.Println(cache.put("pear", "bar").IsSome()) // false

    fmt.Println(cache.get("pear").Unwrap()) // "bar"
    fmt.Println(cache.get("banana").Unwrap()) // "foo"
    fmt.Println(cache.get("apple").IsSome()) // false
}
```

[build badge]: https://github.com/lovelysunlight/lru-go/actions/workflows/ci.yaml/badge.svg
[build status]: https://github.com/lovelysunlight/lru-go/actions/workflows/ci.yaml
[codecov badge]: https://codecov.io/gh/lovelysunlight/lru-go/branch/master/graph/badge.svg
[coverage status]: https://codecov.io/gh/lovelysunlight/lru-go
[license badge]: https://img.shields.io/badge/license-MIT-blue.svg
[license]: https://raw.githubusercontent.com/lovelysunlight/lru-go/master/LICENSE
[LRU Cache implementation in a Rust crate]: https://github.com/jeromefroe/lru-rs
