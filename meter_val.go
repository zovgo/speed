package speed

import (
	"sync"
	"sync/atomic"
	"time"
)

// NewMeterVal creates new speed meter instance.
func NewMeterVal[T any](limit, maxObjects int64) *MeterVal[T] {
	maxObjects = max(0, maxObjects)
	return &MeterVal[T]{
		limit:      limit,
		maxObjects: maxObjects,
		objects:    make([]object[T], 0, maxObjects),
	}
}

// MeterVal is used to count how much object is being created in N duration. It
// will store T object.
type MeterVal[T any] struct {
	sync.Mutex // Allow users to control lock by themselves

	limit, maxObjects int64
	objects           []object[T]
}

// Append appends new object to meter objects.
func (meter *MeterVal[T]) Append(obj T) {
	meter.Lock()
	defer meter.Unlock()
	meter.AppendUnsafe(obj)
}

const DefaultMaxObjects = 100

// AppendUnsafe appends object to speed meter without locking mutex.
func (meter *MeterVal[T]) AppendUnsafe(obj T) {
	if (meter.maxObjects <= 0 && len(meter.objects) >= DefaultMaxObjects) ||
		(meter.maxObjects > 0 && len(meter.objects) >= int(meter.maxObjects)) {
		// Maximum limit exceeded, limiting objects amount
		meter.objects = meter.objects[1:]
	}
	meter.objects = append(meter.objects, newObject(obj))
}

// LimitExceeded returns true if limit of objects was exceeded by n duration.
func (meter *MeterVal[T]) LimitExceeded(n time.Duration) bool {
	limit := meter.Limit()
	if limit <= 0 {
		return false
	}
	return len(meter.Objects(n)) >= limit
}

// Objects returns all objects, that are added not longer than n ago. If user
// will enter 0 as n, function will return all objects.
func (meter *MeterVal[T]) Objects(n time.Duration) []T {
	meter.Lock()
	defer meter.Unlock()
	return meter.ObjectsUnsafe(n)
}

// ObjectsUnsafe ...
func (meter *MeterVal[T]) ObjectsUnsafe(n time.Duration) (objects []T) {
	for _, obj := range meter.objects {
		if n == 0 || time.Since(obj.createdAt) <= n {
			objects = append(objects, obj.val)
		}
	}
	return
}

// Limit returns the speed meter limit.
func (meter *MeterVal[T]) Limit() int {
	return int(atomic.LoadInt64(&meter.limit))
}

// SetLimit updates speed meter limit.
func (meter *MeterVal[T]) SetLimit(n int) {
	atomic.StoreInt64(&meter.limit, int64(n))
}
