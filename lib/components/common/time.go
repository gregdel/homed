package common

import (
	"encoding/json"
	"sync"
	"time"
)

var zeroTime time.Time

// Time represents a concurrency-safe time.
type Time struct {
	mu sync.RWMutex
	v  time.Time
}

// Load returns a *time.Time
func (t *Time) Load() *time.Time {
	if t == nil {
		return nil
	}

	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.v.IsZero() {
		return nil
	}

	v := t.v
	return &v
}

// Store sets the time from a *time.Time
func (t *Time) Store(v *time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if v != nil {
		t.v = *v
		return
	}

	t.v = zeroTime
}

// MarshalJSON implements the json.Marshaler interface
func (t *Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Load())
}

// UnmarshalJSON implements the json.Marshaler interface
func (t *Time) UnmarshalJSON(data []byte) error {
	var v *time.Time
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	if v == nil || v.IsZero() {
		t.Store(nil)
	} else {
		t.Store(v)
	}

	return nil
}
