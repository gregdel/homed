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

	h.log = logger.With(zap.String("friendly_name", h.Name))

	ticker := time.NewTicker(30 * time.Second)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			h.updateTemperatureMode()
			h.setTRVTarget()
		case <-h.Events.Incoming:
			// We're only subcribed to sensors, let's not check the event ID
			h.updateTemperature()
		}

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
	}

	h.Current.Store(value)
}

func (h *HomedTemperature) updateTemperatureMode() {
	log := h.log

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

func (h *HomedTemperature) setTRVTarget() {
	log := h.log

	target, err := h.TemperatureTarget()
	if err != nil {
		log.Warn("failed to get temperature target", zap.Error(err))
		return
	}

	roomTemperature := h.Current.Load()

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
