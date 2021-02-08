package homedtemperature

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/components/common"
	"github.com/prometheus/client_golang/prometheus"
)

// Make sure that the module is a temperature controller
var _ components.TemperatureControllerInternal = (*HomedTemperature)(nil)

func init() {
	components.Register(components.TypeHomedTemperature, New)
}

// Data represents the data of HomedTemperature
type Data struct {
	Current      float64                    `json:"current"`
	Target       float64                    `json:"target"`
	Mode         components.TemperatureMode `json:"mode"`
	ManualTarget float64                    `json:"manual_target"`
	ManualUntil  *time.Time                 `json:"manual_until,omitempty"`
}

// HomedTemperature is a component that handles temperatures
type HomedTemperature struct {
	common.ScheduledComponent

	Data
}

// New returns a new temperature component
func New() components.Component {
	sc := common.NewScheduledComponent()
	return &HomedTemperature{
		ScheduledComponent: *sc,
		Data: Data{
			Mode: components.TemperatureModeAuto,
		},
	}
}

// CurrentTarget returns the current target according to the mode
func (h *HomedTemperature) CurrentTarget() float64 {
	if h.Mode == components.TemperatureModeAuto {
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
		Mode           components.TemperatureMode `json:"mode"`
		ManualTarget   float64                    `json:"manual_target"`
		ManualUntil    *time.Time                 `json:"manual_until,omitempty"`
		ManualDuration string                     `json:"manual_duration"`
	}{}

	if err := json.Unmarshal(cmd, &data); err != nil {
		return err
	}

	switch data.Mode {
	case components.TemperatureModeAuto:
		data.ManualTarget = h.Target
		data.ManualUntil = nil
	case components.TemperatureModeFixed:
		data.ManualUntil = nil
	case components.TemperatureModeDuration:
		d, err := time.ParseDuration(data.ManualDuration)
		if err != nil {
			return err
		}
		t := time.Now().Add(d)
		data.ManualUntil = &t
	case components.TemperatureModeUntilDate:
		if data.ManualUntil == nil {
			return fmt.Errorf("components: homed_temperature: missing date")
		}
	case components.TemperatureModeNextTimeBlock:
		data.ManualUntil = nil
	default:
		return nil
	}

	if data.ManualUntil != nil && time.Now().After(*data.ManualUntil) {
		return fmt.Errorf("components: homed_temperature: date is in the past")
	}

	h.ManualUntil = data.ManualUntil
	h.ManualTarget = data.ManualTarget
	h.Mode = data.Mode

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

// Temperature implements the TemperatureController interface
func (h *HomedTemperature) Temperature() (float64, error) {
	return h.Current, nil
}

// SetTemperature implements the TemperatureSetter interface
func (h *HomedTemperature) SetTemperature(temperature float64) error {
	h.Current = temperature
	return h.PublishState()
}

// TemperatureTarget implements the TemperatureController interface
func (h *HomedTemperature) TemperatureTarget() (float64, error) {
	if h.Mode == components.TemperatureModeAuto {
		return h.Target, nil
	}

	return h.ManualTarget, nil
}

// SetTemperatureTarget implements the TemperatureController interface
func (h *HomedTemperature) SetTemperatureTarget(target float64) error {
	h.Target = target
	return h.PublishState()
}

// TemperatureMode implements the TemperatureController interface
func (h *HomedTemperature) TemperatureMode() (components.TemperatureMode, error) {
	return h.Mode, nil
}

// SetTemperatureMode implements the TemperatureController interface
func (h *HomedTemperature) SetTemperatureMode(mode components.TemperatureMode) error {
	h.Mode = mode
	return h.PublishState()
}

// TemperatureManualTarget implements the TemperatureController interface
func (h *HomedTemperature) TemperatureManualTarget() (float64, error) {
	return h.ManualTarget, nil
}

// SetTemperatureManualTarget implements the TemperatureController interface
func (h *HomedTemperature) SetTemperatureManualTarget(target float64) error {
	h.ManualTarget = target
	return h.PublishState()
}

// SetTemperatureModeManualUntil implements the TemperatureControllerInternal interface
func (h *HomedTemperature) SetTemperatureModeManualUntil(until *time.Time) error {
	h.ManualUntil = until
	return h.PublishState()
}

// TemperatureModeManualUntil implements the TemperatureControllerInternal interface
func (h *HomedTemperature) TemperatureModeManualUntil() (*time.Time, error) {
	return h.ManualUntil, nil
}
