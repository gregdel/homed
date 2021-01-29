package common

// BinarySensor represents a generic sensor
type BinarySensor struct {
	Component

	On bool `json:"on"`
}

// IsOn implements the BinarySensor interface
func (bs *BinarySensor) IsOn() bool {
	return bs.On
}
