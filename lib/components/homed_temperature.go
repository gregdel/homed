package components

import (
	"encoding/json"
	"fmt"
	"time"

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
	HomedTemperatureModeAuto          HomedTemperatureMode = "auto"
	HomedTemperatureModeFixed         HomedTemperatureMode = "fixed"
	HomedTemperatureModeDuration      HomedTemperatureMode = "duration"
	HomedTemperatureModeUntilDate     HomedTemperatureMode = "until_date"
	HomedTemperatureModeNextTimeBlock HomedTemperatureMode = "next_time_block"
)

// HomedTemperatureData represents the data of HomedTemperature
type HomedTemperatureData struct {
	Current      float64              `json:"current"`
	Target       float64              `json:"target"`
	Mode         HomedTemperatureMode `json:"mode"`
	ManualTarget float64              `json:"manual_target"`
	ManualUntil  *time.Time           `json:"manual_until,omitempty"`
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

// CurrentTarget returns the current target according to the mode
func (s *HomedTemperature) CurrentTarget() float64 {
	if s.Mode == HomedTemperatureModeAuto {
		return s.Target
	}

	return s.ManualTarget
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

// ExecCommand implements the Component interface
func (s *HomedTemperature) ExecCommand(client mqtt.Client, cmd []byte) error {
	data := struct {
		Mode           HomedTemperatureMode `json:"mode"`
		ManualTarget   float64              `json:"manual_target"`
		ManualUntil    *time.Time           `json:"manual_until,omitempty"`
		ManualDuration string               `json:"manual_duration"`
	}{}

	if err := json.Unmarshal(cmd, &data); err != nil {
		return err
	}

	switch data.Mode {
	case HomedTemperatureModeAuto:
		data.ManualTarget = s.Target
		data.ManualUntil = nil
	case HomedTemperatureModeFixed:
		data.ManualUntil = nil
	case HomedTemperatureModeDuration:
		d, err := time.ParseDuration(data.ManualDuration)
		if err != nil {
			return err
		}
		t := time.Now().Add(d)
		data.ManualUntil = &t
	case HomedTemperatureModeUntilDate:
		if data.ManualUntil == nil {
			return fmt.Errorf("components: homed_temperature: missing date")
		}
	case HomedTemperatureModeNextTimeBlock:
		data.ManualUntil = nil
	default:
		return nil
	}

	if data.ManualUntil != nil && time.Now().After(*data.ManualUntil) {
		return fmt.Errorf("components: homed_temperature: date is in the past")
	}

	s.Mode = data.Mode
	s.ManualTarget = data.ManualTarget
	s.ManualUntil = data.ManualUntil

	return s.PublishState(client)
}

// PublishState publishes the mqtt state of the component
func (s *HomedTemperature) PublishState(client mqtt.Client) error {
	// Publish the current state
	data, err := json.Marshal(s.HomedTemperatureData)
	if err != nil {
		return err
	}

	token := client.Publish(s.stateTopic, 0, true, data)
	token.Wait()
	return token.Error()
}

// Update implements the Component interface
func (s *HomedTemperature) Update(value []byte) error {
	return json.Unmarshal(value, s)
}
