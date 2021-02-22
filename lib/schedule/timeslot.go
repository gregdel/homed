package schedule

import "fmt"

// TimeSlot holds a time slot
type TimeSlot struct {
	ID string `json:"id"`
	// Required
	Start Time `json:"start"`
	// Optional
	Stop  *Time   `json:"stop"`
	Value float64 `json:"value"`
}

// String implements the stringer interface
func (t *TimeSlot) String() string {
	if t.Stop == nil {
		return fmt.Sprintf("%f from %s", t.Value, t.Start)
	}

	return fmt.Sprintf("%f from %s to %s", t.Value, t.Start, t.Stop)
}

func (t *TimeSlot) generateID() {
	if t.ID != "" {
		return
	}

	t.ID = generateID(t)
	return
}
