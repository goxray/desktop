package connlist

import (
	"time"
)

type NetworkRecorder interface {
	// Start should start recording data.
	Start()
	// Read should return values for uplink for each previous RecordInterval.
	// Number of values returned must match Written.
	Read() []float64
	// Written should return values for downlink for each previous RecordInterval.
	// Number of values returned must match Written.
	Written() []float64
	// BytesRead should return the total number of bytes for uplink.
	BytesRead() int
	// BytesWritten should return the total number of bytes for downlink.
	BytesWritten() int
	RecordInterval() time.Duration
}

func (c *Item) Read() []float64 {
	return c.recorder.Read()
}

func (c *Item) Written() []float64 {
	return c.recorder.Written()
}

func (c *Item) BytesRead() int {
	return c.recorder.BytesRead()
}

func (c *Item) BytesWritten() int {
	return c.recorder.BytesWritten()
}

func (c *Item) RecordInterval() time.Duration {
	return c.recorder.RecordInterval()
}
