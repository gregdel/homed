package homed

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
	return 0
	// // For now, we only return the first value of the component type "Temperature"
	// var temperature float64
	// var controlled bool
	// var tuyaTemp float64

	// for _, device := range r.Devices {
	// 	for _, component := range device.Components {
	// 		switch component.Type() {
	// 		case components.TypeHomedTemperature:
	// 			controlled = true
	// 		case components.TypeTemperature:
	// 			c := component.(*components.Temperature)
	// 			temperature = c.Value
	// 		case components.TypeTuyaTRV:
	// 			c := component.(*components.TuyaTRV)
	// 			tuyaTemp = c.Temperature
	// 		}

	// 		if controlled && temperature != 0 {
	// 			return temperature
	// 		}
	// 	}
	// }

	// return tuyaTemp
}
