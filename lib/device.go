package homed

// Device represents a device
type Device struct {
	Room *Room  `json:"-"`
	Name string `json:"name"`
}

// NewDevice creates a new device
func NewDevice(name string) *Device {
	return &Device{
		Name: name,
	}
}
