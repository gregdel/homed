package homed

import (
	"encoding/json"

	"github.com/prometheus/client_golang/prometheus"
)

// SensorRTL433 is a sensor that handles rtl_433 signals
type SensorRTL433 struct {
	BaseSensor

	BoilerState bool
}

// NewSensorRTL433 returns a new humidity sensor
func NewSensorRTL433() *SensorRTL433 {
	return &SensorRTL433{}
}

// Type implements the Sensor interface
func (s *SensorRTL433) Type() SensorType {
	return SensorTypeRTL433
}

// Init implements the Sensor interface
func (s *SensorRTL433) Init() error {
	if s.device == nil {
		return ErrMissingDevice
	}

	c := prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "homed_boiler_state",
		},
		func() float64 {
			if s.BoilerState {
				return 1
			}
			return 0
		},
	)

	return prometheus.Register(c)
}

// Update implements the Sensor interface
func (s *SensorRTL433) Update(value []byte) error {
	data := struct {
		Model string `json:"model"`
		Rows  []struct {
			Len  int    `json:"len"`
			Data string `json:"data"`
		} `json:"rows"`
	}{}

	if err := json.Unmarshal(value, &data); err != nil {
		return nil
	}

	if data.Model != "boiler" {
		return nil
	}

	if len(data.Rows) != 3 {
		return nil
	}

	code := data.Rows[2].Data
	switch code {
	case "ff67b7efb7fbb8":
		s.BoilerState = true
	case "ff67b7ffdbfddf":
		s.BoilerState = true
	case "ff67b7ffdbfedf":
		s.BoilerState = false
	case "ff67b7efb7fdb8":
		s.BoilerState = false
	}

	return nil
}
