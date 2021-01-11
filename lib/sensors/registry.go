package sensors

import "fmt"

var registeredSensors map[string]func() Sensor

func register(t Type, fn func() Sensor) {
	if registeredSensors == nil {
		registeredSensors = map[string]func() Sensor{}
	}

	name := string(t)
	_, ok := registeredSensors[name]
	if ok {
		err := fmt.Errorf("sensors: sensor %s already regitered", name)
		panic(err)
	}

	registeredSensors[name] = fn
}

// newSensor returns a new sensor from a type name
func newSensor(typeName string) (Sensor, error) {
	fn, ok := registeredSensors[typeName]
	if !ok {
		return nil, fmt.Errorf("sensors: sensor %s is not registered", typeName)
	}

	return fn(), nil
}
