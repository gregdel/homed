package common

import (
	"fmt"
	"math"
	"strconv"
)

// RollerShutter represents a generic roller shutter
type RollerShutter struct {
	GenericSensor
}

// OpenedAt implements the RollerShutter interface
func (rs *RollerShutter) OpenedAt() float64 {
	return rs.Value
}

// IsOpen implements the RollerShutter interface
func (rs *RollerShutter) IsOpen() bool {
	if rs.Value == 100 {
		return true
	}
	return false
}

// IsClosed implements the RollerShutter interface
func (rs *RollerShutter) IsClosed() bool {
	return !rs.IsOpen()
}

// OpenAt implements the RollerShutter interface
func (rs *RollerShutter) OpenAt(v float64) error {
	value := int(math.Round(v))

	if value < 0 || value > 100 {
		return fmt.Errorf("common: invalid roller shutter opening value: %d", value)
	}

	return rs.WriteCommand([]byte(strconv.Itoa(value)))
}
