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
func (r *Room) Temperature(cs []components.Component) float64 {
	var temperature float64
	var found float64 = 0

	for _, c := range cs {
		if c.Type() == components.TypeHomedTemperature {
			continue
		}

		if c.Type() == components.TypeTuyaTRV {
			// Don't use this for now
			continue
		}

		tc, ok := c.(components.TemperatureGetter)
		if !ok {
			continue
		}

		found = found + 1
		t, _ := tc.Temperature()
		temperature = (temperature + t) / found
	}

	return temperature
}
