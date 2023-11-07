package homedtemperature

import (
	"context"
	"errors"
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

	log := h.LoggerWithFields(logger)
	h.log = log

	ticker := time.NewTicker(30 * time.Second)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			h.updateTemperatureMode()
		case event := <-h.Events.Incoming:
			h.mu.Lock()
			_, ok := h.sensors[event.ID]
			h.mu.Unlock()
			if ok {
				h.updateTemperature()
				log.Debug("updating temperature")
			}
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
		if c.Type() == components.TypeZigbeeTRV {
			h.mu.Lock()
			h.trvs[c.ID()] = c.(components.TemperatureController)
			h.mu.Unlock()
		}

		sensor, ok := c.(components.TemperatureGetter)
		if ok {
			if c.Type() == components.TypeHomedTemperature ||
				c.Type() == components.TypeZigbeeTRV {
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

func (h *HomedTemperature) updateTemperature() error {
	var temperature float64
	var found float64

	for _, sensor := range h.sensors {
		t, err := sensor.Temperature()
		if err != nil {
			return err
		}
		temperature += t
		found++
	}

	h.SetTemperature(temperature / found)
	return nil
}

func (h *HomedTemperature) updateTemperatureMode() {
	log := h.log

	// Unset manual until
	manualUntil, _ := h.TemperatureModeManualUntil()
	if manualUntil != nil && time.Now().After(*manualUntil) {
		h.Mode.Store(string(components.TemperatureModeAuto))

		h.mu.Lock()
		h.ManualUntil = nil
		h.mu.Unlock()
	}

	// Get the current schedule value
	schedule := h.Schedule()
	scheduledTarget, opportunistic := schedule.Values()
	h.Target.Store(scheduledTarget)

	mode, _ := h.TemperatureMode()

	// Only the mode auto can be opportunistic
	if mode != components.TemperatureModeAuto {
		opportunistic = false
	}
	h.Opportunistic.Store(opportunistic)

	// Transform until next change to manual until
	if mode == components.TemperatureModeUntilNextChange {
		nextChange := schedule.NextChange()
		if nextChange != nil {
			if err := h.SetTemperatureModeManualUntil(nextChange); err != nil {
				log.Warn("failed to set the temperature manual mode until", zap.Error(err))
				return
			}
		}
	}

	// Set the isHeating property
	isHeating := false
	current, _ := h.Temperature()
	currentTarget, _ := h.TemperatureTarget()
	if (current < currentTarget) && !opportunistic {
		isHeating = true
	}
	if current == 0 {
		log.Warn("the temperature is reported to be 0, this is unlikely, let's ignore this for now")
		isHeating = false
	}
	if currentTarget == 0 {
		log.Warn("the temperature target is 0, this is unlikely, let's ignore this for now")
		isHeating = false
	}
	h.Heating.Store(isHeating)

	if err := h.PublishState(); err != nil {
		log.Warn("failed to publish state", zap.Error(err))
		return
	}
}

func (h *HomedTemperature) setTRVTarget(log *zap.Logger) {
	target, err := h.TemperatureTarget()
	if err != nil {
		log.Warn("failed to get temperature target", zap.Error(err))
		return
	}

	roomTemperature, _ := h.Temperature()

	for _, trv := range h.trvs {
		trvLogger := trv.LoggerWithFields(log)

		// Compute the new calibration
		trvTemperature, err := trv.Temperature()
		if err != nil {
			trvLogger.Warn("failed to get trv temperature", zap.Error(err))
			continue
		}

		if h.Params.CalibrateTRV {
			if err := h.recalibrateTRV(log, trv, roomTemperature, trvTemperature); err != nil {
				trvLogger.Warn("failed to calibrate the trv", zap.Error(err))
				continue
			}
		}

		trvTarget, err := trv.TemperatureTarget()
		if err != nil {
			trvLogger.Warn("failed to get trv's temperature target", zap.Error(err))
			continue
		}

		if trvTarget == target {
			continue
		}

		trvLogger.Info("should configure the trv to the target",
			zap.Float64("target", target))

		err = trv.SetTemperatureTarget(target)
		if err != nil {
			trvLogger.Warn("failed to set trv's temperature target", zap.Error(err))
			continue
		}
	}
}

func (h *HomedTemperature) recalibrateTRV(logger *zap.Logger, trv components.TemperatureController, roomTemperature, trvTemperature float64) error {
	log := trv.LoggerWithFields(logger)

	calibration, err := trv.TemperatureCalibration()
	if err != nil {
		if errors.Is(err, components.ErrOperatingInProgress) {
			log.Debug("calibration already in progress")
			return nil
		}

		return err
	}

	trvMesuredTemperature := trvTemperature - calibration

	delta := roomTemperature - trvMesuredTemperature

	// Only keep on decimal of precision
	delta = math.Round(delta*10) / 10

	if math.Abs(delta) > h.Params.CalibrationMaxOffset {
		log.Debug("invalid calibration, ignoring for now", zap.Float64("calibration", delta))
		return nil
	}

	if math.Abs(calibration-delta) > h.Params.CalibrationThreshold {
		log.Info("recalibrating the trv",
			zap.Float64("old_calibration", calibration),
			zap.Float64("new_calibration", delta),
			zap.Float64("calibration_diff", math.Abs(calibration-delta)),
			zap.Float64("calibration_threshold",
				h.Params.CalibrationThreshold),
		)
		err = trv.SetTemperatureCalibration(delta)
		if err != nil {
			if errors.Is(err, components.ErrOperatingInProgress) {
				log.Debug("operation already in progress")
			} else {
				log.Warn("failed to calibrate trv", zap.Error(err))
			}
		}
	}

	return nil
}
