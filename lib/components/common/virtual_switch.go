package common

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/gregdel/homed/lib/components"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	components.Register(components.TypeVirtualSwitch, NewVirtualSwitch)
}

// NewVirtualSwitch returns a new virtual button
func NewVirtualSwitch() components.Component {
	return &VirtualSwitch{}
}

// VirtualSwitch represents a virtual button that can trigger multiple switches
type VirtualSwitch struct {
	BinarySensor
	counterMu sync.RWMutex
	counter   float64

	log      *slog.Logger
	switches []components.Switch
}

type VirtualSwitchSnapshot struct {
	ID           string             `json:"id"`
	UpdatedAt    *time.Time         `json:"updated_at"`
	FriendlyName string             `json:"friendly_name"`
	Hide         bool               `json:"hide"`
	Device       *components.Device `json:"device"`
	On           bool               `json:"on"`
	Counter      float64            `json:"counter"`
}

// Type implements the Component interface
func (v *VirtualSwitch) Type() components.Type {
	return components.TypeVirtualSwitch
}

// Update implements the Component interface
func (v *VirtualSwitch) Update(payload []byte) error {
	if string(payload) != "single" {
		v.log.Debug("update with invalid payload, skipping",
			slog.String("payload", string(payload)))
		return nil
	}

	v.log.Debug("update called", slog.Int("switches", len(v.switches)))
	if err := v.Toggle(); err != nil {
		return err
	}
	v.incrementCounter()
	return nil
}

// TurnOn implements the switch interface
func (v *VirtualSwitch) TurnOn() error {
	for _, sw := range v.switches {
		if err := sw.TurnOn(); err != nil {
			return err
		}
	}

	v.SetOn(true)
	return v.PostUpdate()
}

// TurnOff implements the switch interface
func (v *VirtualSwitch) TurnOff() error {
	for _, sw := range v.switches {
		if err := sw.TurnOff(); err != nil {
			return err
		}
	}

	v.SetOn(false)
	return v.PostUpdate()
}

// Toggle implements the Switch interface
func (v *VirtualSwitch) Toggle() error {
	if v.IsOn() {
		return v.TurnOff()
	}

	return v.TurnOn()
}

func (v *VirtualSwitch) updateState() {
	isOn := false
	for _, sw := range v.switches {
		if sw.IsOn() {
			isOn = true
			break
		}
	}

	if isOn == v.IsOn() {
		return
	}

	v.log.Debug("updating state", slog.Bool("new_state", isOn))
	v.SetOn(isOn)
	if err := v.PostUpdate(); err != nil {
		v.log.Error("failed to update state", slog.Any("error", err))
	}
}

// WriteCommand implements the Component interface
func (v *VirtualSwitch) WriteCommand(data []byte) error {
	// Called by the http app
	return v.Toggle()
}

// Run implements the Component interface
func (v *VirtualSwitch) Run(ctx context.Context, logger *slog.Logger,
	inventory *components.Components) error {
	v.Events.Incoming = components.NewEventChannel()

	v.log = v.LoggerWithFields(logger)

	params := struct {
		Switches []string `yaml:"switches"`
	}{Switches: []string{}}

	if err := v.YAMLParams.Decode(&params); err != nil {
		return err
	}

	for _, id := range params.Switches {
		swComponent, err := inventory.Get(id)
		if err != nil {
			v.log.Error("failed to get switch", slog.Any("error", err))
		}

		sw, ok := (swComponent).(components.Switch)
		if !ok {
			return fmt.Errorf("component %q is not a switch", id)
		}

		swComponent.Subscribe(v.ID(), v.Events.Incoming)
		v.switches = append(v.switches, sw)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-v.Events.Incoming:
			v.updateState()
		}
	}
}

// Collectors implements the Component interface
func (v *VirtualSwitch) Collectors(labels prometheus.Labels) []prometheus.Collector {
	return []prometheus.Collector{
		components.GaugeCollector("switch", labels, func() float64 {
			if v.IsOn() {
				return 1
			}

			return 0
		}),
		components.CounterCollector("virtual_switch_counter", labels,
			func() float64 { return v.CounterValue() },
		),
	}
}

func (v *VirtualSwitch) CounterValue() float64 {
	v.counterMu.RLock()
	defer v.counterMu.RUnlock()

	return v.counter
}

func (v *VirtualSwitch) incrementCounter() {
	v.counterMu.Lock()
	defer v.counterMu.Unlock()

	v.counter++
}

func (v *VirtualSwitch) Snapshot() any {
	return VirtualSwitchSnapshot{
		ID:           v.ID(),
		UpdatedAt:    v.UpdatedAt.Load(),
		FriendlyName: v.FriendlyName(),
		Hide:         v.Hide,
		Device:       v.Device(),
		On:           v.IsOn(),
		Counter:      v.CounterValue(),
	}
}
