package sensors

import (
	"encoding/json"

	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	register(TypeRTL433, NewRTL433)
}

// RTL433 is a sensor that handles rtl_433 signals
type RTL433 struct {
	BoilerState bool `json:"boiler_state"`
}

// NewRTL433 returns a new humidity sensor
func NewRTL433() Sensor {
	return &RTL433{}
}

// Type implements the Sensor interface
func (s *RTL433) Type() Type {
	return TypeRTL433
}

// Collectors implements the Sensor interface
func (s *RTL433) Collectors(labels prometheus.Labels) []prometheus.Collector {
	return []prometheus.Collector{
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{Name: "homed_boiler_state"},
			func() float64 {
				if s.BoilerState {
					return 1
				}
				return 0
			},
		),
	}
}

// Update implements the Sensor interface
func (s *RTL433) Update(value []byte) error {
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
