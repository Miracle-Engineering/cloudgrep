package provider

import (
	"sync/atomic"
	"time"
)

// FakeProviderResourceConfig stores configuration on a specific type for this fake provider
type FakeProviderResourceConfig struct {
	// Count is the number of resource returned by the fetch function
	Count int

	// ErrorAfterCount sets the error that the fetch function returns after sending the `Count` resources
	ErrorAfterCount error

	// DelayBefore sets a sleep duration before writing to the output channel
	DelayBefore time.Duration

	runCount int32
}

// RunCount returns the number of completed runs the fetch function has made, even if it errors out.
// This is useful for determining when the fetch function has returned
func (c *FakeProviderResourceConfig) RunCount() int {
	return int(atomic.LoadInt32(&c.runCount))
}

func (c *FakeProviderResourceConfig) addRun() {
	atomic.AddInt32(&c.runCount, 1)
}
