package homedtemperature

import (
	"context"
	"math"
	"time"

	"github.com/gregdel/homed/lib/components"
	"go.uber.org/zap"
)

// Run implements the Component interface
func (h *HomedTemperature) Run(ctx context.Context, logger *zap.Logger, inventory *components.Components) error {
	if err := h.setup(inventory); err != nil {
		return err
	}

	h.log = h.LoggerWithFields(logger)

	ticker := time.NewTicker(30 * time.Second)

	h.updateTemperatureMode()
	h.updateTemperature()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			h.updateTemperatureMode()
		case <-h.Events.Incoming:
			// We're only subcribed to sensors, let's not check the event ID
		}

		h.updateTemperature()
		h.handleBinaryTRV()
		if err := h.PublishState(); err != nil {
			h.log.Warn("failed to publish state", zap.Error(err))
		}
	}
}

func (h *HomedTemperature) setup(inventory *components.Components) error {
	h.Events.Incoming = make(chan components.Event)

	if err := h.YAMLParams.Decode(&h.Params); err != nil {
		return err
	}

	room := h.Device().Room

	list := inventory.ListByRoom(room)
	for _, c := range list {
		if c.Type() == components.TypeBinaryTRV {
			h.mu.Lock()
			h.binTRVs[c.ID()] = c.(components.Switch)
			h.mu.Unlock()
		}

		sensor, ok := c.(components.TemperatureGetter)
		if ok {
			if c.Type() == components.TypeHomedTemperature {
				continue
			}

			h.mu.Lock()
			h.sensors[sensor.ID()] = sensor
			h.mu.Unlock()

			sensor.Subscribe(h.ID(), h.Events.Incoming)
		}
	}

	return nil
}

func (h *HomedTemperature) updateTemperature() {
	var temperature float64
	var found float64

	for _, sensor := range h.sensors {
		t, err := sensor.Temperature()
		if err != nil {
			h.log.Error("failed to get temperature", zap.Error(err))
			continue
		}
		temperature += t
		found++
	}

	var value float64
	if found > 0 {
		value = (temperature / found)
	}

	if math.IsNaN(value) {
		value = 0
	} else {
		value = math.Round(value*10) / 10
	}

	h.Current.Store(value)
}

func (h *HomedTemperature) updateTemperatureMode() {
	log := h.log

	if !h.IsOn() {
		return
	}

	// Unset manual until
	manualUntil := h.ManualUntil.Load()
	if manualUntil != nil && time.Now().After(*manualUntil) {
		h.Mode.Store(string(components.TemperatureModeAuto))
		h.ManualUntil.Store(nil)
	}

	// Get the current schedule value
	schedule := h.Schedule()
	scheduledTarget, opportunistic := schedule.Values()
	h.Target.Store(scheduledTarget)

	mode := components.TemperatureMode(h.Mode.Load())

	// Only the mode auto can be opportunistic
	if mode != components.TemperatureModeAuto {
		opportunistic = false
	}
	h.Opportunistic.Store(opportunistic)

	// Transform until next change to manual until
	if mode == components.TemperatureModeUntilNextChange {
		nextChange := schedule.NextChange()
		if nextChange != nil {
			h.ManualUntil.Store(nextChange)
		}
	}

	// Set the isHeating property
	isHeating := false
	current := h.Current.Load()
	currentTarget, _ := h.TemperatureTarget()
	if (current < currentTarget) && !opportunistic {
		isHeating = true
	}
	if current == 0 {
		log.Warn("the temperature is reported to be 0, this is unlikely, let's ignore this for now")
		isHeating = false
		return
	}
	if currentTarget == 0 {
		log.Warn("the temperature target is 0, this is unlikely, let's ignore this for now")
		isHeating = false
		return
	}

	h.Heating.Store(isHeating)
}
