package homedtemperature

import (
	"context"
	"log/slog"
	"math"
	"time"

	"github.com/gregdel/homed/lib/components"
)

// Run implements the Component interface
func (h *HomedTemperature) Run(ctx context.Context, logger *slog.Logger, inventory *components.Components) error {
	if err := h.setup(inventory); err != nil {
		return err
	}

	h.log = h.LoggerWithFields(logger)

	// Wait for the component to be updated by the retained data.
	try := 0
	tryWait := 3 * time.Second
	ticker := time.NewTicker(tryWait)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-h.Events.Incoming:
			h.log.Info("homed temperature ignoring events for now")
		case <-ticker.C:
		}

		if h.hasCurrentTemperature() {
			h.log.Info("homed temperature updated from retained data")
			break
		}

		h.log.Info("homed temperature not ready yet")
		ticker.Reset(tryWait)
		try++
		if try >= 10 {
			h.log.Warn("homed temperature doesn't have retained data")
			break
		}
	}

	ticker.Reset(30 * time.Second)

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
			h.log.Warn("failed to publish state", slog.Any("error", err))
		}
	}
}

func (h *HomedTemperature) setup(inventory *components.Components) error {
	h.Events.Incoming = components.NewEventChannel()

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
			h.log.Error("failed to get temperature", slog.Any("error", err))
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

	h.mu.Lock()
	h.Current = value
	h.mu.Unlock()
}

func (h *HomedTemperature) updateTemperatureMode() {
	log := h.log

	if !h.IsOn() {
		return
	}

	schedule := h.Schedule()
	scheduledTarget, opportunistic := schedule.Values()

	h.mu.RLock()
	mode := components.TemperatureMode(h.Mode)
	h.mu.RUnlock()

	var nextChange *time.Time
	if mode == components.TemperatureModeUntilNextChange {
		nextChange = schedule.NextChange()
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	h.updateTemperatureModeLocked(time.Now(), scheduledTarget, opportunistic, nextChange, log)
}

func (h *HomedTemperature) updateTemperatureModeLocked(now time.Time, scheduledTarget float64, opportunistic bool, nextChange *time.Time, log *slog.Logger) {
	if !h.On {
		return
	}

	// Unset manual until
	if h.ManualUntil != nil && now.After(*h.ManualUntil) {
		h.Mode = string(components.TemperatureModeAuto)
		h.setManualUntilLocked(nil)
	}

	// Get the current schedule value
	h.Target = scheduledTarget

	mode := components.TemperatureMode(h.Mode)

	// Only the mode auto can be opportunistic
	if mode != components.TemperatureModeAuto {
		opportunistic = false
	}
	h.Opportunistic = opportunistic

	// Transform until next change to manual until
	if mode == components.TemperatureModeUntilNextChange && nextChange != nil {
		h.setManualUntilLocked(nextChange)
	}

	// Set the isHeating property
	isHeating := false
	current := h.Current
	currentTarget := h.temperatureTargetLocked()
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

	h.Heating = isHeating
}

func (h *HomedTemperature) hasCurrentTemperature() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return h.Current > 0
}
