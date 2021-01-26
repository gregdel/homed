package base

import (
	"math"
	"strconv"
)

// Float64Component represents a basic component
type Float64Component struct {
	Value float64 `json:"value"`
}

// Update implements the Component interface
func (fs *Float64Component) Update(value []byte) error {
	v, err := strconv.ParseFloat(string(value), 64)
	if err != nil {
		return err
	}

	// TODO: find a better solution to handle NaN values
	if math.IsNaN(v) {
		v = -99999
	}

	fs.Value = v
	return nil
}
