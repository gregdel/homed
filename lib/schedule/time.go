package schedule

import "fmt"

// Time represents time
type Time struct {
	Hour   int `json:"hour"`
	Minute int `json:"minute"`
	Second int `json:"second"`
}

// NewTime returns a new time
func NewTime(hour, minute, second int) Time {
	return Time{
		Hour:   hour,
		Minute: minute,
		Second: second,
	}
}

// NewTimePointer returns a new time pointer
func NewTimePointer(hour, minute, second int) *Time {
	t := NewTime(hour, minute, second)
	return &t
}

// String returns the time in a printable format
func (t Time) String() string {
	return fmt.Sprintf("%02d:%02d:%02d", t.Hour, t.Minute, t.Second)
}

// Before reports whether the time instant t is before u.
func (t *Time) Before(u Time) bool {
	if t.Hour < u.Hour {
		return true
	}

	if t.Hour > u.Hour {
		return false
	}

	// Same hour

	if t.Minute < u.Minute {
		return true
	}

	if t.Minute > u.Minute {
		return false
	}

	// Same minute

	return t.Second < u.Second
}

// After reports whether the time instant t is after u.
func (t *Time) After(u Time) bool {
	return !t.Before(u)
}
