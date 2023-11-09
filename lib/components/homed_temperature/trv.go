package homedtemperature

import (
	"errors"
	"math"

	"github.com/gregdel/homed/lib/components"
	"go.uber.org/zap"
)

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
