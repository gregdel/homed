package homed

// ActionType is a type representing the actions
type ActionType string

// Action types
var (
	ActionTypeSetTemperature ActionType = "set_temperature"
)

// Action represents an actions
type Action interface {
	Set(interface{}) error
}

type baseAction struct {
	mqttTopic string
}

// Set implements the Action interface
func (ba *baseAction) Set(input interface{}) error {
	return nil
}

// ActionSetTemperature is an action to set the temperature
type ActionSetTemperature struct {
	baseAction
}

// Set implements the Action interface
func (a *ActionSetTemperature) Set(input interface{}) error {
	// Ensure that the given data is a number between 5 and 35

	return nil
}

// NewActionSetTemperature returns a new ActionTypeSetTemperature
func NewActionSetTemperature() *ActionSetTemperature {
	return &ActionSetTemperature{}
}
