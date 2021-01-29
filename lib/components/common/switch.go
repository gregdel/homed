package common

// Switch represents a generic switch
type Switch struct {
	Component
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
