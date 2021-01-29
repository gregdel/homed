package common

import (
	"math"
	"strconv"
)

// GenericSensor represents a generic sensor
type GenericSensor struct {
	Component

	Value float64 `json:"value"`
}

// Update implements the Component interface
func (gs *GenericSensor) Update(value []byte) error {
	v, err := strconv.ParseFloat(string(value), 64)
	if err != nil {
		return err
	}

	// TODO: find a better solution to handle NaN values
	if math.IsNaN(v) {
		v = -99999
	}

	gs.Value = v
	return nil
}

// SensorValue implements the Sensor interface
func (gs *GenericSensor) SensorValue() float64 {
	return gs.Value
}
