package common

import (
	"bytes"
	"fmt"
)

// Switch represents a generic switch
type Switch struct {
	BinarySensor

	payloadOn  []byte
	payloadOff []byte
}

// NewSwitch returns a new switch
func NewSwitch(payloadOn, payloadOff []byte) Switch {
	return Switch{
		payloadOn:  payloadOn,
		payloadOff: payloadOff,
	}
}

// SetOn implements the Status interface
func (s *Switch) SetOn() error {
	return s.WriteCommand(s.payloadOn)
}

// SetOff implements the Status interface
func (s *Switch) SetOff() error {
	return s.WriteCommand(s.payloadOff)
}

// Set implements the Switch interface
func (s *Switch) Set(state bool) error {
	if state {
		return s.SetOn()
	}

	return s.SetOff()
}

// Toggle implements the Switch interface
func (s *Switch) Toggle() error {
	if s.IsOn() {
		return s.SetOff()
	}

	return s.SetOn()
}

// Update implements the Component interface
func (s *Switch) Update(value []byte) error {
	switch {
	case bytes.Equal(value, s.payloadOn):
		s.On = true
	case bytes.Equal(value, s.payloadOff):
		s.On = false
	default:
		return fmt.Errorf("switch: invalid payload: %s", value)
	}

	return nil
}
