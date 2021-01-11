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
