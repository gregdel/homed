package schedule

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ErrInvalidTime represents an invalid time error
var ErrInvalidTime = errors.New("schedule: invalid time")

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

// MarshalJSON implements the marshaler interface
func (t *Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// UnmarshalJSON implements the marshaler interface
func (t *Time) UnmarshalJSON(data []byte) error {
	var input string
	if err := json.Unmarshal(data, &input); err != nil {
		return err
	}

	parts := strings.Split(input, ":")
	if len(parts) != 2 && len(parts) != 3 {
		return ErrInvalidTime
	}

	type conv struct {
		p *string
		t *int
	}

	convs := []conv{
		{p: &parts[0], t: &t.Hour},
		{p: &parts[1], t: &t.Minute},
	}
	if len(parts) == 3 {
		convs = append(convs, conv{p: &parts[2], t: &t.Second})
	}

	for _, s := range convs {
		v, err := strconv.Atoi(*s.p)
		if err != nil {
			return ErrInvalidTime
		}

		*s.t = v
	}

	return nil
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
