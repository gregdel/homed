package components

import (
	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	register(TypeHomedTemperature, NewHomedTemperature)
}

// HomedTemperatureMode represents a temperature control mode
type HomedTemperatureMode string

// Available modes
var (
	HomedTemperatureModeAuto   HomedTemperatureMode = "auto"
	HomedTemperatureModeManual HomedTemperatureMode = "manual"
)

// HomedTemperatureData represents the data of HomedTemperature
type HomedTemperatureData struct {
	Current float64              `json:"current"`
	Target  float64              `json:"target"`
	Mode    HomedTemperatureMode `json:"mode"`
}

// HomedTemperature is a component that handles temperatures
type HomedTemperature struct {
	baseComponent
	HomedTemperatureData
}

// NewHomedTemperature returns a new temperature component
func NewHomedTemperature() Component {
	return &HomedTemperature{
		HomedTemperatureData: HomedTemperatureData{
			Mode: HomedTemperatureModeAuto,
		},
	}
}

// Type implements the Component interface
func (s *HomedTemperature) Type() Type {
	return TypeHomedTemperature
}

// Collectors implements the Component interface
func (s *HomedTemperature) Collectors(labels prometheus.Labels) []prometheus.Collector {
	prefix := "homed_temperature_control_"
	return []prometheus.Collector{
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name:        prefix + "current",
				ConstLabels: labels,
			},
			func() float64 { return s.Current },
		),
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name:        prefix + "target",
				ConstLabels: labels,
			},
			func() float64 { return s.Target },
		),
	}
}

// Update implements the Component interface
func (s *HomedTemperature) Update(value []byte) error {
	return json.Unmarshal(value, s)
}

// ExecCommand implements the Component interface
func (s *HomedTemperature) ExecCommand(client mqtt.Client, cmd []byte) error {
	// Update the internal state
	if err := json.Unmarshal(cmd, s); err != nil {
		return err
	}

	// Publish the current state
	data, err := json.Marshal(s.HomedTemperatureData)
	if err != nil {
		return err
	}

	token := client.Publish(s.stateTopic, 0, true, data)
	token.Wait()
	return token.Error()
}
