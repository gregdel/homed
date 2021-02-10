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
	// For now, we only return the first value of the component type "Temperature"
	var temperature float64
	var controlled bool
	var tuyaTemp float64

	for _, c := range cs {
		if c.Type() == components.TypeHomedTemperature {
			controlled = true
			continue
		}

		tc, ok := c.(components.TemperatureGetter)
		if !ok {
			continue
		}

		if c.Type() == components.TypeTuyaTRV {
			// TODO: handle the error
			tuyaTemp, _ = tc.Temperature()
			continue
		}

		temperature, _ = tc.Temperature()
	}

	if controlled && temperature != 0 {
		return temperature
	}

	return tuyaTemp
}
