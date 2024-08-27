package list

type Interface[K comparable, V any] interface {
	Root() *Entry[K, V]
}
