package speed

import (
	"sync"
	"time"
)

// NewMeter creates new Meter instance.
func NewMeter(limit, maxObjects int64) *Meter {
	meter := NewMeterVal[struct{}](limit, maxObjects)
	return &Meter{
		Locker: &meter.Mutex,
		meter:  meter,
	}
}

// Meter is speed meter that is not accepts any value.
type Meter struct {
	Locker *sync.Mutex
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

// LimitExceeded ...
func (meter *Meter) LimitExceeded(n time.Duration) bool {
	return meter.meter.LimitExceeded(n)
}
