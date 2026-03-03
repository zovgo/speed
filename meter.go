package speed

import (
	"sync"
	"time"
)

func NewMeter(limit, maxObjects int64) *Meter {
	meter := NewMeterVal[struct{}](limit, maxObjects)
	return &Meter{
		Locker: &meter.RWMutex,
		meter:  meter,
	}
}

type Meter struct {
	Locker *sync.RWMutex
	meter  *MeterVal[struct{}]
}

var zeroStruct struct{}

// Append ...
func (meter *Meter) Append() {
	meter.meter.Append(zeroStruct)
}

// AppendUnsafe ...
func (meter *Meter) AppendUnsafe() {
	meter.meter.AppendUnsafe(zeroStruct)
}

// Objects ...
func (meter *Meter) Objects(n time.Duration) int {
	return len(meter.meter.Objects(n))
}

// ObjectsUnsafe ...
func (meter *Meter) ObjectsUnsafe(n time.Duration) (objects int) {
	return len(meter.meter.ObjectsUnsafe(n))
}

// ObjectsLen ...
func (meter *Meter) ObjectsLen(n time.Duration) int {
	return meter.meter.ObjectsLen(n)
}

// ObjectsLenUnsafe ...
func (meter *Meter) ObjectsLenUnsafe(n time.Duration) (objects int) {
	return meter.meter.ObjectsLenUnsafe(n)
}

// LimitExceeded ...
func (meter *Meter) LimitExceeded(n time.Duration) bool {
	return meter.meter.LimitExceeded(n)
}

// LimitExceededUnsafe ...
func (meter *Meter) LimitExceededUnsafe(n time.Duration) bool {
	return meter.meter.LimitExceededUnsafe(n)
}
