package sensors

import (
	"strconv"
)

// float64Sensor represents a basic sensor
type float64Sensor struct {
	Value float64 `json:"value"`
}

// Update implements the Sensor interface
func (fs *float64Sensor) Update(value []byte) error {
	v, err := strconv.ParseFloat(string(value), 64)
	if err != nil {
		return err
	}
	fs.Value = v
	return nil
}
