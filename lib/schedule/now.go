package schedule

import (
	"time"
)

// Function to be overwitten by the tests
var now func() time.Time

func init() {
	// Set the now function to time.Now()
	setNow(nil)
}

// Helper to set the current time to a fake time
func setNow(t *Time) {
	if t == nil {
		now = time.Now
	} else {
		n := time.Now()
		ft := time.Date(
			n.Year(),
			n.Month(),
			n.Day(),
			t.Hour,
			t.Minute,
			t.Second,
			0,
			n.Location(),
		)
		now = func() time.Time {
			return ft
		}
	}
}
