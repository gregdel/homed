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
	ManualUntil   common.Time    `json:"manual_until,omitempty"`
	Heating       atomic.Bool    `json:"heating"`
	Opportunistic atomic.Bool    `json:"opportunistic"`
	On            atomic.Bool    `json:"on"`
}

// HomedTemperature is a component that handles temperatures
type HomedTemperature struct {
	common.ScheduledComponent
	Params config.TemperatureControl
	log    *zap.Logger

	sensors map[string]components.TemperatureGetter
	trvs    map[string]components.TemperatureController
	binTRVs map[string]components.Switch

	mu sync.RWMutex
	Data
}

// New returns a new temperature component
func New() components.Component {
	h := &HomedTemperature{
		sensors: map[string]components.TemperatureGetter{},
		trvs:    map[string]components.TemperatureController{},
		binTRVs: map[string]components.Switch{},
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
		components.GaugeCollector("temperature_control_on", labels,
			func() float64 {
				if h.IsOn() {
					return 1
				}

				return 0
			},
		),
	}
}

// ExecCommand implements the Component interface
func (h *HomedTemperature) ExecCommand(cmd []byte) error {
	data := struct {
		On             bool                       `json:"on"`
		Mode           components.TemperatureMode `json:"mode"`
		ManualTarget   float64                    `json:"manual_target"`
		ManualUntil    *time.Time                 `json:"manual_until,omitempty"`
		ManualDuration string                     `json:"manual_duration"`
	}{}

	if err := json.Unmarshal(cmd, &data); err != nil {
		return err
	}

	if data.ManualUntil != nil && time.Now().After(*data.ManualUntil) {
		return fmt.Errorf("components: homed_temperature: date is in the past")
	}

	switch data.Mode {
	case components.TemperatureModeOnOff:
		h.On.Store(data.On)
	case components.TemperatureModeAuto:
		h.ManualUntil.Store(nil)
	case components.TemperatureModeFixed:
		h.ManualUntil.Store(nil)
		h.ManualTarget.Store(data.ManualTarget)
	case components.TemperatureModeDuration:
		d, err := time.ParseDuration(data.ManualDuration)
		if err != nil {
			return err
		}
		t := time.Now().Add(d)
		h.ManualUntil.Store(&t)
		h.ManualTarget.Store(data.ManualTarget)
	case components.TemperatureModeUntilDate:
		if data.ManualUntil == nil {
			return fmt.Errorf("components: homed_temperature: missing date")
		}
		h.ManualUntil.Store(data.ManualUntil)
		h.ManualTarget.Store(data.ManualTarget)
	case components.TemperatureModeUntilNextChange:
		h.ManualUntil.Store(nil)
		h.ManualTarget.Store(data.ManualTarget)
	default:
		return nil
	}

	if data.Mode != components.TemperatureModeOnOff {
		h.Mode.Store(string(data.Mode))
	}

	h.updateTemperatureMode()
	return h.PublishState()
}

// PublishState publishes the mqtt state of the component
func (h *HomedTemperature) PublishState() error {
	data, err := json.Marshal(&h.Data)
	if err != nil {
		return err
	}

	return h.Component.ExecCommand(data)
}

// Update implements the Component interface
func (h *HomedTemperature) Update(value []byte) error {
	return json.Unmarshal(value, &h.Data)
}

// TemperatureTarget implements the TemperatureController interface
func (h *HomedTemperature) TemperatureTarget() (float64, error) {
	if h.Mode.Load() == string(components.TemperatureModeAuto) {
		return h.Target.Load(), nil
	}

	return h.ManualTarget.Load(), nil
}

// IsHeating implements the TemperatureController interface
func (h *HomedTemperature) IsHeating() bool {
	return h.IsOn() && h.Heating.Load()
}

// IsOn implements the Switch interface
func (h *HomedTemperature) IsOn() bool {
	return h.On.Load()
}

// TurnOn implements the Switch interface
func (h *HomedTemperature) TurnOn() error {
	h.On.Store(true)
	return h.PublishState()
}

// TurnOff implements the Switch interface
func (h *HomedTemperature) TurnOff() error {
	h.On.Store(false)
	return h.PublishState()
}

// Toggle implements the Switch interface
func (h *HomedTemperature) Toggle() error {
	if h.IsOn() {
		return h.TurnOff()
	}

	return h.TurnOn()
}
