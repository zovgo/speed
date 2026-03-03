package speed

import (
	"sync"
	"time"
)

func NewMeterVal[T any](limit, maxObjects int64) *MeterVal[T] {
	maxObjects = min(DefaultMaxObjects, max(0, maxObjects))
	return &MeterVal[T]{
		limit:      limit,
		maxObjects: maxObjects,
		objects:    make([]object[T], 0, maxObjects),
	}
}

type MeterVal[T any] struct {
	sync.RWMutex

	limit, maxObjects int64
	objects           []object[T]
}

func (meter *MeterVal[T]) Append(obj T) {
	meter.Lock()
	defer meter.Unlock()
	meter.AppendUnsafe(obj)
}

const DefaultMaxObjects int64 = 100

func (meter *MeterVal[T]) AppendUnsafe(obj T) {
	if meter.maxObjects <= 0 {
		meter.maxObjects = DefaultMaxObjects
	}
	if meter.maxObjects > 0 && len(meter.objects) >= int(meter.maxObjects) {
		meter.objects = meter.objects[1:]
	}
	meter.objects = append(meter.objects, newObject(obj))
}

func (meter *MeterVal[T]) LimitExceeded(n time.Duration) bool {
	meter.RLock()
	defer meter.RUnlock()
	return meter.LimitExceededUnsafe(n)
}

func (meter *MeterVal[T]) LimitExceededUnsafe(n time.Duration) bool {
	limit := meter.Limit()
	if limit <= 0 {
		return false
	}
	return meter.ObjectsLenUnsafe(n) >= limit
}

func (meter *MeterVal[T]) Objects(n time.Duration) []T {
	meter.RLock()
	defer meter.RUnlock()
	return meter.ObjectsUnsafe(n)
}

func (meter *MeterVal[T]) ObjectsUnsafe(n time.Duration) []T {
	objects := make([]T, 0, len(meter.objects))
	for _, obj := range meter.objects {
		if n > 0 && time.Since(obj.createdAt) > n {
			continue
		}
		objects = append(objects, obj.val)
	}
	return objects
}

func (meter *MeterVal[T]) ObjectsLen(n time.Duration) int {
	meter.RLock()
	defer meter.RUnlock()
	return meter.ObjectsLenUnsafe(n)
}

func (meter *MeterVal[T]) ObjectsLenUnsafe(n time.Duration) (length int) {
	for _, obj := range meter.objects {
		if n > 0 && time.Since(obj.createdAt) > n {
			continue
		}
		length++
	}
	return length
}

func (meter *MeterVal[T]) Limit() int {
	return int(meter.limit)
}
