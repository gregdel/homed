package common

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
	return rs.Value == 100
}

// IsClosed implements the RollerShutter interface
func (rs *RollerShutter) IsClosed() bool {
	return !rs.IsOpen()
}

// Open implements the RollerShutter interface
func (rs *RollerShutter) Open() error {
	return rs.WriteCommand([]byte("open"))
}

// Close implements the RollerShutter interface
func (rs *RollerShutter) Close() error {
	return rs.WriteCommand([]byte("close"))
}

// Stop implements the RollerShutter interface
func (rs *RollerShutter) Stop() error {
	return rs.WriteCommand([]byte("stop"))
}
