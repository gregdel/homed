package common

import (
	"encoding/json"
	"time"

	"go.uber.org/atomic"
)

// Time represents an atomic time
type Time struct {
	atomic.Time
}

// MarshalJSON implements the json.Marshaler interface
func (t *Time) MarshalJSON() ([]byte, error) {
	var v *time.Time = nil

	value := t.Time.Load()
	if !value.IsZero() {
		v = &value
	}

	return json.Marshal(v)
}
