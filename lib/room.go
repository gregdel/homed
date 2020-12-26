package homed

import (
	"fmt"
)

// Room represent a room
type Room struct {
	Name string

	Devices []*Device
}

type roomConfigFormat struct {
	Name        string   `yaml:"name"`
	DeviceNames []string `yaml:"devices"`
}

func (r *Room) configFormat() interface{} {
	return &roomConfigFormat{}
}

func (r *Room) fromConfig(data interface{}, h *Homed) error {
	config, ok := data.(*roomConfigFormat)
	if !ok {
		return ErrInvalidConfigFormat
	}

	r.Name = config.Name

	r.Devices = []*Device{}
	for _, name := range config.DeviceNames {
		file, err := h.Get(FileTypeDevice, name)
		if err != nil {
			return err
		}

		r.Devices = append(r.Devices, file.(*Device))
	}

	return nil
}

func (r *Room) toConfig() (interface{}, error) {
	config := &roomConfigFormat{
		Name:        r.Name,
		DeviceNames: []string{},
	}

	for _, d := range r.Devices {
		config.DeviceNames = append(config.DeviceNames, d.Name)
	}

	return config, nil
}

// FileName implements the File interface
func (r *Room) FileName() string {
	return r.Name
}

// FileType implements the File interface
func (r *Room) FileType() FileType {
	return FileTypeRoom
}

// Pretty implements the Pretty interface
func (r *Room) Pretty() string {
	return fmt.Sprintf("Room name: %s\n", r.Name)
}

// NewRoom returns a new room
func NewRoom(name string) *Room {
	return &Room{Name: name}
}
