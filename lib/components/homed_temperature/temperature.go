package homedtemperature

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gregdel/homed/lib/components"
	base "github.com/gregdel/homed/lib/components/base_component"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	components.Register(components.TypeHomedTemperature, New)
}

// Mode represents a temperature control mode
type Mode string

// Available modes
var (
	ModeAuto          Mode = "auto"
	ModeFixed         Mode = "fixed"
	ModeDuration      Mode = "duration"
	ModeUntilDate     Mode = "until_date"
	ModeNextTimeBlock Mode = "next_time_block"
)

// Data represents the data of HomedTemperature
type Data struct {
	Current      float64    `json:"current"`
	Target       float64    `json:"target"`
	Mode         Mode       `json:"mode"`
	ManualTarget float64    `json:"manual_target"`
	ManualUntil  *time.Time `json:"manual_until,omitempty"`
}

// HomedTemperature is a component that handles temperatures
type HomedTemperature struct {
	base.Component
	Data
}

// New returns a new temperature component
func New() components.Component {
	return &HomedTemperature{
		Data: Data{
			Mode: ModeAuto,
		},
	}
}

// CurrentTarget returns the current target according to the mode
func (h *HomedTemperature) CurrentTarget() float64 {
	if h.Mode == ModeAuto {
		return h.Target
	}

	return h.ManualTarget
}

// Type implements the Component interface
func (h *HomedTemperature) Type() components.Type {
	return components.TypeHomedTemperature
}

// Collectors implements the Component interface
func (h *HomedTemperature) Collectors(labels prometheus.Labels) []prometheus.Collector {
	prefix := "homed_temperature_control_"
	return []prometheus.Collector{
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name:        prefix + "current",
				ConstLabels: labels,
			},
			func() float64 { return h.Current },
		),
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name:        prefix + "target",
				ConstLabels: labels,
			},
			func() float64 { return h.Target },
		),
	}
}

// ExecCommand implements the Component interface
func (h *HomedTemperature) ExecCommand(cmd []byte) error {
	data := struct {
		Mode           Mode       `json:"mode"`
		ManualTarget   float64    `json:"manual_target"`
		ManualUntil    *time.Time `json:"manual_until,omitempty"`
		ManualDuration string     `json:"manual_duration"`
	}{}

	if err := json.Unmarshal(cmd, &data); err != nil {
		return err
	}

	switch data.Mode {
	case ModeAuto:
		data.ManualTarget = h.Target
		data.ManualUntil = nil
	case ModeFixed:
		data.ManualUntil = nil
	case ModeDuration:
		d, err := time.ParseDuration(data.ManualDuration)
		if err != nil {
			return err
		}
		t := time.Now().Add(d)
		data.ManualUntil = &t
	case ModeUntilDate:
		if data.ManualUntil == nil {
			return fmt.Errorf("components: homed_temperature: missing date")
		}
	case ModeNextTimeBlock:
		data.ManualUntil = nil
	default:
		return nil
	}

	if data.ManualUntil != nil && time.Now().After(*data.ManualUntil) {
		return fmt.Errorf("components: homed_temperature: date is in the past")
	}

	h.Mode = data.Mode
	h.ManualTarget = data.ManualTarget
	h.ManualUntil = data.ManualUntil

	return h.PublishState()

}

// PublishState publishes the mqtt state of the component
func (h *HomedTemperature) PublishState() error {
	data, err := json.Marshal(h.Data)
	if err != nil {
		return err
	}

	return h.Component.ExecCommand(data)
}

// Update implements the Component interface
func (h *HomedTemperature) Update(value []byte) error {
	return json.Unmarshal(value, h)
}
