package common

import (
	"context"
	"fmt"

	"github.com/gregdel/homed/lib/components"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/atomic"
	"go.uber.org/zap"
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
	Counter atomic.Float64 `json:"counter"`

	log      *zap.Logger
	switches []components.Switch
}

// Type implements the Component interface
func (v *VirtualSwitch) Type() components.Type {
	return components.TypeVirtualSwitch
}

// Update implements the Component interface
func (v *VirtualSwitch) Update(payload []byte) error {
	if string(payload) != "single" {
		v.log.Debug("update with invalid payload, skipping",
			zap.String("payload", string(payload)))
		return nil
	}

	v.log.Debug("update called", zap.Int("switches", len(v.switches)))
	v.Toggle()
	v.Counter.Add(1)
	return nil
}

// TurnOn implements the switch interface
func (v *VirtualSwitch) TurnOn() error {
	for _, sw := range v.switches {
		sw.TurnOn()
	}

	v.On.Store(true)
	return v.PostUpdate()
}

// TurnOff implements the switch interface
func (v *VirtualSwitch) TurnOff() error {
	for _, sw := range v.switches {
		sw.TurnOff()
	}

	v.On.Store(false)
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

	v.log.Debug("updating state", zap.Bool("new_state", isOn))
	v.On.Store(isOn)
	if err := v.PostUpdate(); err != nil {
		v.log.Error("failed to update state", zap.Error(err))
	}
}

// WriteCommand implements the Component interface
func (v *VirtualSwitch) WriteCommand(data []byte) error {
	// Called by the http app
	v.Toggle()
	return nil
}

// Run implements the Component interface
func (v *VirtualSwitch) Run(ctx context.Context, logger *zap.Logger,
	inventory *components.Components) error {
	v.Events.Incoming = make(chan components.Event)

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
			v.log.Error("failed to get switch", zap.Error(err))
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
			func() float64 { return v.Counter.Load() },
		),
	}
}
