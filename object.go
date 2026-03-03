package speed

import "time"

type object[T any] struct {
	createdAt time.Time
	val       T
}

func newObject[T any](val T) object[T] {
	now := time.Now()
	return object[T]{createdAt: now, val: val}
}
