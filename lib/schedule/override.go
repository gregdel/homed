package schedule

import (
	"fmt"
	"time"
)

// Override represents a time override
type Override struct {
	ID    string    `json:"id"`
	Start time.Time `json:"start"`
	Stop  time.Time `json:"stop"`
	Value float64   `json:"value"`
	On    bool      `json:"on"`
}

// String implements the stringer interface
func (o *Override) String() string {
	return fmt.Sprintf("%f from %s to %s", o.Value, o.Start, o.Stop)
}

// NewOverride returns a new override
func NewOverride(start, stop time.Time, value float64, on bool) *Override {
	o := &Override{
		Start: start,
		Stop:  stop,
		Value: value,
		On:    on,
	}

	o.generateID()
	return o
}

func (o *Override) validate() error {
	if o.Start.IsZero() || o.Stop.IsZero() {
		return ErrInvalidTime
	}

	if o.Stop.Before(o.Start) {
		return ErrStopBeforeStart
	}

	return nil
}

func (o *Override) generateID() {
	if o.ID != "" {
		return
	}

	o.ID = generateID(o)
}
