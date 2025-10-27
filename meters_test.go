package speed_test

import (
	"github.com/go-playground/assert/v2"
	"github.com/zovgo/speed"
	"testing"
	"time"
)

func TestMeter(t *testing.T) {
	const testAmount = 100
	t.Run("default", func(t *testing.T) {
		meter := speed.NewMeter(testAmount+1, testAmount+2)
		start := time.Now()
		for range testAmount {
			meter.Append()
		}
		elapsed := time.Since(start)
		assert.Equal(t, meter.Objects(elapsed), testAmount)
	})
	t.Run("negative amount", func(t *testing.T) {
		meter := speed.NewMeter(-1, -1)
		start := time.Now()
		for range testAmount {
			meter.Append()
		}
		elapsed := time.Since(start)
		assert.Equal(t, meter.Objects(elapsed), testAmount)
	})
	t.Run("limit", func(t *testing.T) {
		const expected = 25
		meter := speed.NewMeter(50, expected)
		start := time.Now()
		for range testAmount {
			meter.Append()
		}
		elapsed := time.Since(start)
		assert.Equal(t, meter.Objects(elapsed), expected)
	})
}
