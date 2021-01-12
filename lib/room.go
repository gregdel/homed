package homed

import "github.com/gregdel/homed/lib/components"

// Room represent a room
type Room struct {
	Name    string    `json:"name"`
	Devices []*Device `json:"devices"`
}

// NewRoom returns a new room
func NewRoom(name string) *Room {
	return &Room{Name: name}
}

// AddDevice adds a device to the room
func (r *Room) AddDevice(device *Device) {
	if r.Devices == nil {
		r.Devices = []*Device{}
	}

	r.Devices = append(r.Devices, device)
}

// Temperature returns the temperature in the room
func (r *Room) Temperature() float64 {
	// For now, we only return the first value of the component type "Temperature"

	for _, device := range r.Devices {
		for _, component := range device.Components {
			if component.Type() == components.TypeTemperature {
				// get the value and publish it to mqtt

				c, ok := component.(*components.Temperature)
				if !ok {
					break
				}

				return c.Value
			}
		}
	}

	return 0
}
