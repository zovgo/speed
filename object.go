package speed

import "time"

// object is used to store objects in speed meter.
type object[T any] struct {
	createdAt time.Time
	val       T
}

// newObject creates new object implementation.
func newObject[T any](val T) object[T] {
	now := time.Now()
	return object[T]{createdAt: now, val: val}
}
