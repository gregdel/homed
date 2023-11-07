package homedtemperature

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/components/common"
	"github.com/gregdel/homed/lib/config"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/atomic"
	"go.uber.org/zap"
)

// Make sure that the module is a temperature controller
var _ components.TemperatureControllerInternal = (*HomedTemperature)(nil)

func init() {
	components.Register(components.TypeHomedTemperature, New)
}

// Data represents the data of HomedTemperature
type Data struct {
	Current       atomic.Float64 `json:"current"`
	Target        atomic.Float64 `json:"target"`
	Mode          atomic.String  `json:"mode"`
	ManualTarget  atomic.Float64 `json:"manual_target"`
	ManualUntil   *time.Time     `json:"manual_until,omitempty"`
	Heating       atomic.Bool    `json:"heating"`
	Opportunistic atomic.Bool    `json:"opportunistic"`
}

// HomedTemperature is a component that handles temperatures
type HomedTemperature struct {
	common.ScheduledComponent
	Params config.TemperatureControl
	log    *zap.Logger

	sensors map[string]components.TemperatureGetter
	trvs    map[string]components.TemperatureController

	mu sync.RWMutex
	Data
}

// New returns a new temperature component
func New() components.Component {
	h := &HomedTemperature{
		sensors: map[string]components.TemperatureGetter{},
		trvs:    map[string]components.TemperatureController{},
	}
	h.Mode.Store(string(components.TemperatureModeAuto))
	h.Heating.Store(false)
	return h
}

// Type implements the Component interface
func (h *HomedTemperature) Type() components.Type {
	return components.TypeHomedTemperature
}

// Collectors implements the Component interface
func (h *HomedTemperature) Collectors(labels prometheus.Labels) []prometheus.Collector {
	return []prometheus.Collector{
		components.GaugeCollector("temperature_control_current", labels,
			func() float64 { return h.Current.Load() },
		),
		components.GaugeCollector("temperature_control_target", labels,
			func() float64 {
				t, err := h.TemperatureTarget()
				if err != nil {
					return 0
				}
				return t
			},
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
		data.ManualTarget = h.Target.Load()
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
	case components.TemperatureModeUntilNextChange:
		data.ManualUntil = nil
	default:
		return nil
	}

	if data.ManualUntil != nil && time.Now().After(*data.ManualUntil) {
		return fmt.Errorf("components: homed_temperature: date is in the past")
	}

	h.mu.Lock()
	h.ManualUntil = data.ManualUntil
	h.mu.Unlock()

	h.ManualTarget.Store(data.ManualTarget)
	h.Mode.Store(string(data.Mode))

	h.updateTemperatureMode()
	return nil
}

// PublishState publishes the mqtt state of the component
func (h *HomedTemperature) PublishState() error {
	h.mu.RLock()
	data, err := json.Marshal(&h.Data)
	h.mu.RUnlock()
	if err != nil {
		return err
	}

	return h.Component.ExecCommand(data)
}

// Update implements the Component interface
func (h *HomedTemperature) Update(value []byte) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	return json.Unmarshal(value, &h.Data)
}

// Temperature implements the TemperatureController interface
func (h *HomedTemperature) Temperature() (float64, error) {
	return h.Current.Load(), nil
}

// SetTemperature implements the TemperatureSetter interface
func (h *HomedTemperature) SetTemperature(temperature float64) error {
	h.Current.Store(temperature)
	return h.PublishState()
}

// TemperatureTarget implements the TemperatureController interface
func (h *HomedTemperature) TemperatureTarget() (float64, error) {
	if h.Mode.Load() == string(components.TemperatureModeAuto) {
		return h.Target.Load(), nil
	}

	return h.ManualTarget.Load(), nil
}

// SetTemperatureTarget implements the TemperatureController interface
func (h *HomedTemperature) SetTemperatureTarget(target float64) error {
	h.Target.Store(target)
	return h.PublishState()
}

// TemperatureMode implements the TemperatureController interface
func (h *HomedTemperature) TemperatureMode() (components.TemperatureMode, error) {
	return components.TemperatureMode(h.Mode.Load()), nil
}

// SetTemperatureMode implements the TemperatureController interface
func (h *HomedTemperature) SetTemperatureMode(mode components.TemperatureMode) error {
	h.Mode.Store(string(mode))
	return h.PublishState()
}

// TemperatureManualTarget implements the TemperatureController interface
func (h *HomedTemperature) TemperatureManualTarget() (float64, error) {
	return h.ManualTarget.Load(), nil
}

// SetTemperatureManualTarget implements the TemperatureController interface
func (h *HomedTemperature) SetTemperatureManualTarget(target float64) error {
	h.ManualTarget.Store(target)
	return h.PublishState()
}

// SetTemperatureModeManualUntil implements the TemperatureControllerInternal interface
func (h *HomedTemperature) SetTemperatureModeManualUntil(until *time.Time) error {
	h.mu.Lock()
	h.ManualUntil = until
	h.mu.Unlock()
	return h.PublishState()
}

// TemperatureModeManualUntil implements the TemperatureControllerInternal interface
func (h *HomedTemperature) TemperatureModeManualUntil() (*time.Time, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.ManualUntil, nil
}

// TemperatureCalibration implements the TemperatureController interface
func (h *HomedTemperature) TemperatureCalibration() (float64, error) {
	return 0, components.ErrNotImplemented
}

// SetTemperatureCalibration implements the TemperatureController interface
func (h *HomedTemperature) SetTemperatureCalibration(c float64) error {
	return components.ErrNotImplemented
}

// IsHeating implements the TemperatureController interface
func (h *HomedTemperature) IsHeating() bool {
	return h.Heating.Load()
}

// SetHeating implements the TemperatureController interface
func (h *HomedTemperature) SetHeating(b bool) error {
	h.Heating.Store(b)
	return h.PublishState()
}
