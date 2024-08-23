package deepcopy

import (
	"reflect"
	"time"
)

// Interface for delegating copy process to type
type Interface[T any] interface {
	DeepCopy() T
}

// Iface is an alias to Copy; this exists for backwards compatibility reasons.
//
//go:inline
func Iface(iface interface{}) interface{} {
	return Copy(iface)
}

// Copy creates a deep copy of whatever is passed to it and returns the copy
// in an interface{}.  The returned value will need to be asserted to the
// correct type.
//
//go:inline
func Copy[T any](src T) T {
	var zero T
	if reflect.DeepEqual(src, zero) {
		return src
	}

	// Make the interface a reflect.Value
	original := reflect.ValueOf(src)

	// If it's a basic type, we don't need to do anything special.
	if willCopy[T](&original) {
		return src
	}

	// Make a copy of the same type as the original.
	cpy := reflect.New(original.Type()).Elem()

	// Recursively copy the original.
	copyRecursive[T](original, cpy)

	// Return the copy as an interface.
	return cpy.Interface().(T)
}

//go:inline
func willCopy[T any](cpy *reflect.Value) bool {
	switch cpy.Type().String() {
	case "bool":
		fallthrough
	case "string":
		fallthrough
	case "float32", "float64":
		fallthrough
	case "int", "int8", "int16", "int32", "int64":
		fallthrough
	case "uint", "uint8", "uint16", "uint32", "uint64":
		return true
	default:
		return false
	}
}

// copyRecursive does the actual copying of the interface. It currently has
// limited support for what it can handle. Add as needed.
func copyRecursive[T any](original, cpy reflect.Value) {
	// check for implement deepcopy.Interface
	if original.CanInterface() {
		if copier, ok := original.Interface().(Interface[T]); ok {
			cpy.Set(reflect.ValueOf(copier.DeepCopy()))
			return
		}
	}

	// handle according to original's Kind
	switch original.Kind() {
	case reflect.Ptr:
		// Get the actual value being pointed to.
		originalValue := original.Elem()

		// if  it isn't valid, return.
		if !originalValue.IsValid() {
			return
		}
		cpy.Set(reflect.New(originalValue.Type()))
		copyRecursive[T](originalValue, cpy.Elem())

	case reflect.Struct:
		t, ok := original.Interface().(time.Time)
		if ok {
			cpy.Set(reflect.ValueOf(t))
			return
		}
		// copy it.
		cpy.Set(original)

	case reflect.Slice:
		if original.IsNil() {
			return
		}
		// Make a new slice and copy each element.
		cpy.Set(reflect.MakeSlice(original.Type(), original.Len(), original.Cap()))
		for i := 0; i < original.Len(); i++ {
			copyRecursive[T](original.Index(i), cpy.Index(i))
		}

	case reflect.Map:
		if original.IsNil() {
			return
		}
		cpy.Set(reflect.MakeMap(original.Type()))
		for _, key := range original.MapKeys() {
			originalValue := original.MapIndex(key)
			copyValue := reflect.New(originalValue.Type()).Elem()
			copyRecursive[T](originalValue, copyValue)
			copyKey := Copy(key.Interface())
			cpy.SetMapIndex(reflect.ValueOf(copyKey), copyValue)
		}
	case reflect.Chan, reflect.Func:
		panic("deepcopy: unsupported `chan` and `function`")
	// case reflect.Interface:
	// 	return
	default:
		cpy.Set(original)
	}
}
