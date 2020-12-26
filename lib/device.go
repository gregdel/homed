package homed

import "fmt"

// Device represents a device
type Device struct {
	Name    string
	Sensors []*Sensor
}

type deviceConfigFormat struct {
	Room        string   `yaml:"room"`
	Name        string   `yaml:"name"`
	SensorNames []string `yaml:"sensors"`
}

func (d *Device) configFormat() interface{} {
	return &deviceConfigFormat{}
}

func (d *Device) fromConfig(data interface{}, h *Homed) error {
	config, ok := data.(*deviceConfigFormat)
	if !ok {
		return ErrInvalidConfigFormat
	}

	d.Name = config.Name

	// TODO: handle sensors
	d.Sensors = []*Sensor{}
	for _, name := range config.SensorNames {
		file, err := h.Get(FileTypeSensor, name)
		if err != nil {
			return err
		}

		// Sensors are generic, we need to duplicate it to use
		// it properly
		sensor := file.(*Sensor)
		sensor = sensor.Duplicate()
		sensor.Device = d

		d.Sensors = append(d.Sensors, sensor)
	}

	return nil
}

func (d *Device) toConfig() (interface{}, error) {
	config := &deviceConfigFormat{
		Name:        d.Name,
		SensorNames: []string{},
	}
	// if d.Room != nil {
	// 	config.Room = d.Room.Name
	// }

	for _, s := range d.Sensors {
		config.SensorNames = append(config.SensorNames, s.Name)
	}

	return config, nil
}

// FileName implements the File interface
func (d *Device) FileName() string {
	return d.Name
}

// FileType implements the File interface
func (d *Device) FileType() FileType {
	return FileTypeDevice
}

// Pretty implements the Pretty interface
func (d *Device) Pretty() string {
	return fmt.Sprintf("Device name: %s\n", d.Name)
}

// NewDevice creates a new device
func NewDevice(name string) *Device {
	return &Device{Name: name}
}
