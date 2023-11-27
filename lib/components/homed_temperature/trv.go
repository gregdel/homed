package homedtemperature

import (
	"errors"
	"math"

	"github.com/gregdel/homed/lib/components"
	"go.uber.org/zap"
)

func (h *HomedTemperature) handleBinaryTRV() {
	if len(h.binTRVs) == 0 || !h.On.Load() {
		return
	}

	current := h.Current.Load()
	target, err := h.TemperatureTarget()
	if err != nil {
		h.log.Warn("failed to get temperature target", zap.Error(err))
		return
	}

	for _, trv := range h.binTRVs {
		if current < target {
			trv.TurnOn()
		} else {
			trv.TurnOff()
		}
	}
}

func (h *HomedTemperature) setTRVTarget() {
	if len(h.trvs) == 0 || !h.On.Load() {
		return
	}

	log := h.log

	target, err := h.TemperatureTarget()
	if err != nil {
		log.Warn("failed to get temperature target", zap.Error(err))
		return
	}

	for _, trv := range h.trvs {
		trvLogger := h.log.With(zap.String("trv", trv.FriendlyName()))

		if h.Params.CalibrateTRV {
			if err := h.recalibrateTRV(trv); err != nil {
				trvLogger.Warn("failed to calibrate the trv", zap.Error(err))
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

func (h *HomedTemperature) recalibrateTRV(trv components.TemperatureController) error {
	log := h.log.With(zap.String("trv", trv.FriendlyName()))

	calibration, err := trv.TemperatureCalibration()
	if err != nil {
		if errors.Is(err, components.ErrOperatingInProgress) {
			log.Debug("calibration already in progress")
			return nil
		}

		return err
	}

	trvTemperature, err := trv.Temperature()
	if err != nil {
		return err
	}

	trvMesuredTemperature := trvTemperature - calibration

	delta := h.Current.Load() - trvMesuredTemperature

	// Only keep one decimal of precision
	delta = math.Round(delta*10) / 10

	if math.Abs(delta) > h.Params.CalibrationMaxOffset {
		log.Debug("invalid calibration, ignoring for now", zap.Float64("calibration", delta))
		return nil
	}

	diff := math.Abs(calibration - delta)
	diff = math.Round(diff*10) / 10
	if diff > h.Params.CalibrationThreshold {
		log.Info("recalibrating the trv",
			zap.Float64("old_calibration", calibration),
			zap.Float64("new_calibration", delta),
			zap.Float64("calibration_diff", diff),
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
