//go:generate mockgen -destination=mocks/recorder_mocks.go -source=recorder.go -package=mocks -typed

package netchart

import (
	"context"
	"slices"
	"sync"
	"time"
)

const bytesToMb = 125000

type Source interface {
	BytesRead() int
	BytesWritten() int
}

type Recorder struct {
	base     Source
	interval time.Duration
	mu       sync.RWMutex

	stopRecording   func()
	done            chan struct{}
	recordedRead    []float64
	recordedWritten []float64
	recordLimit     int
	totalRead       int
	totalWrite      int
}

// NewRecorder creates a default Recorder
// TODO: Decrease detalization for old data as the time goes on to allow for longer ranges charts.
func NewRecorder(s Source) *Recorder {
	return &Recorder{
		base:        s,
		interval:    time.Second, // store data value per interval
		recordLimit: 60 * 2,      // store and record only last 2 minutes of data
		done:        make(chan struct{}),
	}
}

func (r *Recorder) Read() []float64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.recordedRead
}

func (r *Recorder) Written() []float64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.recordedWritten
}

func (r *Recorder) RecordInterval() time.Duration {
	return r.interval
}

func (r *Recorder) Start() {
	var ctx context.Context
	ctx, r.stopRecording = context.WithCancel(context.Background())
	go func() {
		for {
			select {
			case <-ctx.Done():
				r.done <- struct{}{}
				return
			case <-time.After(r.interval):
				func() {
					r.mu.RLock()
					defer r.mu.RUnlock()

					rawRead := float64(r.ReadSinceLast() / bytesToMb)
					rawWritten := float64(r.WrittenSinceLast() / bytesToMb)

					if len(r.recordedRead) > r.recordLimit {
						r.recordedRead = slices.Delete(r.recordedRead, 0, 1)
					}
					r.recordedRead = append(r.recordedRead, rawRead)

					if len(r.recordedWritten) > r.recordLimit {
						r.recordedWritten = slices.Delete(r.recordedWritten, 0, 1)
					}
					r.recordedWritten = append(r.recordedWritten, rawWritten)
				}()
			}
		}

	}()
}

func (r *Recorder) Stop() {
	r.stopRecording()
	<-r.done
}

func (r *Recorder) BytesRead() int {
	return r.totalRead
}

func (r *Recorder) BytesWritten() int {
	return r.totalWrite
}

// ReadSinceLast returns bytes read from last call (upload).
func (r *Recorder) ReadSinceLast() int {
	if r.base == nil {
		return 0
	}
	readSinceLast := r.base.BytesRead() - r.totalRead
	r.totalRead = r.base.BytesRead()

	return max(readSinceLast, 0)
}

// WrittenSinceLast returns bytes written from last call (download).
func (r *Recorder) WrittenSinceLast() int {
	if r.base == nil {
		return 0
	}
	writtenSinceLast := r.base.BytesWritten() - r.totalWrite
	r.totalWrite = r.base.BytesWritten()

	return max(writtenSinceLast, 0)
}
