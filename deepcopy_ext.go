package lru

import (
	"github.com/lovelysunlight/lru-go/internal/deepcopy"
)

type deepCopyExt[K comparable, V any] struct {
	copy bool
}

//go:inline
func (c *deepCopyExt[K, V]) OptionalCopyKey(data K) K {
	if c.copy {
		return deepcopy.Copy(data)
	}
	return data
}

//go:inline
func (c *deepCopyExt[K, V]) OptionalCopyValue(data V) V {
	if c.copy {
		return deepcopy.Copy(data)
	}
	return data
}

//go:inline
func (c *deepCopyExt[K, V]) OptionalCopyKeyN(data []K) []K {
	if c.copy {
		return deepcopy.Copy(data)
	}
	return data
}

//go:inline
func (c *deepCopyExt[K, V]) OptionalCopyValueN(data []V) []V {
	if c.copy {
		return deepcopy.Copy(data)
	}
	return data
}
