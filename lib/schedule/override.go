package schedule

import (
	"time"

	"github.com/google/uuid"
)

// Override represents a time override
type Override struct {
	ID    string    `json:"id"`
	Start time.Time `json:"start"`
	Stop  time.Time `json:"stop"`
	Value float64   `json:"value"`
}

// NewOverride returns a new override
func NewOverride(start, stop time.Time, value float64) *Override {
	o := &Override{
		Start: start,
		Stop:  stop,
		Value: value,
	}
	o.generateID()
	return o
}

func (o *Override) validate() error {
	if o.Start.IsZero() || o.Stop.IsZero() {
		return ErrInvalidTime
	}

	return nil
}

func (o *Override) generateID() {
	if o.ID != "" {
		return
	}

	// Generate a random uuid
	uuid, err := uuid.NewRandom()
	if err != nil {
		return
	}

	// Add a random ID to the timeslot
	o.ID = uuid.String()
}
