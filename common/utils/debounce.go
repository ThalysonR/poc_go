package utils

import (
	"sync"
	"time"
)

type debouncer struct {
	duration time.Duration
	mu       sync.Mutex
	timer    *time.Timer
}

func NewDebouncer(duration time.Duration) func(f func()) {
	d := &debouncer{
		duration: duration,
	}

	return func(f func()) {
		d.debounce(f)
	}
}

func (d *debouncer) debounce(f func()) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		d.timer.Stop()
	}
	d.timer = time.AfterFunc(d.duration, f)
}
