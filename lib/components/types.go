package components

// Switch is an interface to reprensents a switch
type Switch interface {
	IsOn() bool
	Toggle() error
	SetOn() error
	SetOff() error
	Set(state bool) error
}

// Publisher is an interface to publish the component state
type Publisher interface {
	PublishState() error
}
