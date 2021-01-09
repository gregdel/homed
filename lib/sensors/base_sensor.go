package sensors

import "github.com/prometheus/client_golang/prometheus"

// BaseSensor represents a basic sensor
type BaseSensor struct {
	mqttTopic string
}

// Type implements the Sensor interface
func (s *BaseSensor) Type() Type {
	return TypeUnknown
}

// Init implements the Sensor interface
func (s *BaseSensor) Init() error {
	return nil
}

// Update implements the Sensor interface
func (s *BaseSensor) Update(_ []byte) error {
	return nil
}

// Collectors implements the Sensor interface
func (s *BaseSensor) Collectors() []prometheus.Collector {
	return nil
}
