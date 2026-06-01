package homedtemperature

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/components/common"
	"github.com/gregdel/homed/lib/config"
	"github.com/prometheus/client_golang/prometheus"
)

// Make sure that the module is a temperature controller
var _ components.TemperatureControllerInternal = (*HomedTemperature)(nil)

func init() {
	components.Register(components.TypeHomedTemperature, New)
}

// Data represents the data of HomedTemperature
type Data struct {
	Current       float64    `json:"current"`
	Target        float64    `json:"target"`
	Mode          string     `json:"mode"`
	ManualTarget  float64    `json:"manual_target"`
	ManualUntil   *time.Time `json:"manual_until"`
	Heating       bool       `json:"heating"`
	Opportunistic bool       `json:"opportunistic"`
	On            bool       `json:"on"`
}

// StateSnapshot represents the MQTT state JSON for HomedTemperature.
type StateSnapshot struct {
	Current       float64    `json:"current"`
	Target        float64    `json:"target"`
	Mode          string     `json:"mode"`
	ManualTarget  float64    `json:"manual_target"`
	ManualUntil   *time.Time `json:"manual_until"`
	Heating       bool       `json:"heating"`
	Opportunistic bool       `json:"opportunistic"`
	On            bool       `json:"on"`
}

// Snapshot represents the HTTP and websocket JSON values for HomedTemperature.
type Snapshot struct {
	ID           string             `json:"id"`
	UpdatedAt    *time.Time         `json:"updated_at"`
	FriendlyName string             `json:"friendly_name"`
	Hide         bool               `json:"hide"`
	Device       *components.Device `json:"device"`
	StateSnapshot
}

// HomedTemperature is a component that handles temperatures
type HomedTemperature struct {
	common.ScheduledComponent
	Params config.TemperatureControl
	log    *slog.Logger

	sensors map[string]components.TemperatureGetter
	binTRVs map[string]components.Switch

	mu sync.RWMutex
	Data
}

// New returns a new temperature component
func New() components.Component {
	h := &HomedTemperature{
		sensors: map[string]components.TemperatureGetter{},
		binTRVs: map[string]components.Switch{},
	}
	h.Mode = string(components.TemperatureModeAuto)
	h.Heating = false
	return h
}

// StateSnapshot returns the current MQTT state JSON shape.
func (h *HomedTemperature) StateSnapshot() StateSnapshot {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return h.stateSnapshotLocked()
}

func (h *HomedTemperature) stateSnapshotLocked() StateSnapshot {
	return StateSnapshot{
		Current:       h.Current,
		Target:        h.Target,
		Mode:          h.Mode,
		ManualTarget:  h.ManualTarget,
		ManualUntil:   copyTimePtr(h.ManualUntil),
		Heating:       h.Heating,
		Opportunistic: h.Opportunistic,
		On:            h.On,
	}
}

// Snapshot returns the current HTTP and websocket JSON values shape.
func (h *HomedTemperature) Snapshot() any {
	return Snapshot{
		ID:            h.ID(),
		UpdatedAt:     h.UpdatedAt.Load(),
		FriendlyName:  h.FriendlyName(),
		Hide:          h.Hide,
		Device:        h.Device(),
		StateSnapshot: h.StateSnapshot(),
	}
}

func copyTimePtr(t *time.Time) *time.Time {
	if t == nil {
		return nil
	}

	copied := *t
	return &copied
}

func (h *HomedTemperature) setManualUntilLocked(t *time.Time) {
	h.ManualUntil = copyTimePtr(t)
}

func (h *HomedTemperature) temperatureTargetLocked() float64 {
	if h.Mode == string(components.TemperatureModeAuto) {
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
	return []prometheus.Collector{
		components.GaugeCollector("temperature_control_current", labels,
			func() float64 {
				h.mu.RLock()
				defer h.mu.RUnlock()
				return h.Current
			},
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

	h.mu.Lock()
	switch data.Mode {
	case components.TemperatureModeOnOff:
		h.On = data.On
	case components.TemperatureModeAuto:
		h.setManualUntilLocked(nil)
	case components.TemperatureModeFixed:
		h.setManualUntilLocked(nil)
		h.ManualTarget = data.ManualTarget
	case components.TemperatureModeDuration:
		d, err := time.ParseDuration(data.ManualDuration)
		if err != nil {
			h.mu.Unlock()
			return err
		}
		t := time.Now().Add(d)
		h.setManualUntilLocked(&t)
		h.ManualTarget = data.ManualTarget
	case components.TemperatureModeUntilDate:
		if data.ManualUntil == nil {
			h.mu.Unlock()
			return fmt.Errorf("components: homed_temperature: missing date")
		}
		h.setManualUntilLocked(data.ManualUntil)
		h.ManualTarget = data.ManualTarget
	case components.TemperatureModeUntilNextChange:
		h.setManualUntilLocked(nil)
		h.ManualTarget = data.ManualTarget
	default:
		h.mu.Unlock()
		return nil
	}

	if data.Mode != components.TemperatureModeOnOff {
		h.Mode = string(data.Mode)
	}
	h.mu.Unlock()

	h.updateTemperatureMode()
	return h.PublishState()
}

// PublishState publishes the mqtt state of the component
func (h *HomedTemperature) PublishState() error {
	data, err := json.Marshal(h.StateSnapshot())
	if err != nil {
		return err
	}

	return h.Component.ExecCommand(data)
}

// Update implements the Component interface
func (h *HomedTemperature) Update(value []byte) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	data := h.Data
	if data.ManualUntil != nil {
		data.ManualUntil = copyTimePtr(data.ManualUntil)
	}
	if err := json.Unmarshal(value, &data); err != nil {
		return err
	}

	data.ManualUntil = copyTimePtr(data.ManualUntil)
	h.Data = data
	return nil
}

// TemperatureTarget implements the TemperatureController interface
func (h *HomedTemperature) TemperatureTarget() (float64, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return h.temperatureTargetLocked(), nil
}

// IsHeating implements the TemperatureController interface
func (h *HomedTemperature) IsHeating() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return h.On && h.Heating
}

// IsOn implements the Switch interface
func (h *HomedTemperature) IsOn() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return h.On
}

// TurnOn implements the Switch interface
func (h *HomedTemperature) TurnOn() error {
	h.mu.Lock()
	h.On = true
	h.mu.Unlock()

	return h.PublishState()
}

// TurnOff implements the Switch interface
func (h *HomedTemperature) TurnOff() error {
	h.mu.Lock()
	h.On = false
	h.mu.Unlock()

	return h.PublishState()
}

// Toggle implements the Switch interface
func (h *HomedTemperature) Toggle() error {
	if h.IsOn() {
		return h.TurnOff()
	}

	return h.TurnOn()
}
