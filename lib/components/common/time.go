package common

import (
	"encoding/json"
	"time"

	"go.uber.org/atomic"
)

var zeroTime time.Time

// Time represents an atomic time
type Time struct {
	v atomic.Time
}

// Load returns a *time.Time
func (t *Time) Load() *time.Time {
	v := t.v.Load()
	if v.IsZero() {
		return nil
	}

	return &v
}

// Store sets the time from a *time.Time
func (t *Time) Store(v *time.Time) {
	if v != nil {
		t.v.Store(*v)
		return
	}

	t.v.Store(zeroTime)
}

// MarshalJSON implements the json.Marshaler interface
func (t Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Load())
}

// UnmarshalJSON implements the json.Marshaler interface
func (t *Time) UnmarshalJSON(data []byte) error {
	var v time.Time
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	if v.IsZero() {
		t.Store(nil)
	} else {
		t.Store(&v)
	}

	return nil
}
