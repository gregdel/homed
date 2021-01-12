package components

import (
	"strconv"
)

// float64Component represents a basic component
type float64Component struct {
	Value float64 `json:"value"`
}

// Update implements the Component interface
func (fs *float64Component) Update(value []byte) error {
	v, err := strconv.ParseFloat(string(value), 64)
	if err != nil {
		return err
	}
	fs.Value = v
	return nil
}
