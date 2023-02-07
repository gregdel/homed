package tempd

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/gregdel/homed/lib/apps"
	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/config"
	"go.uber.org/zap"
)

func init() {
	apps.Register(app())
}

type tempd struct {
	config config.TemperatureControl

	logger     *zap.Logger
	components *components.Components

	boiler components.Switch

	// Wether we're in a rising or falling phase of the hysteresis algorithm
	rising map[string]bool

	rooms map[string]components.TemperatureControllerInternal
	trvs  map[string][]components.TemperatureController
}

func app() *tempd {
	return &tempd{
		rising: map[string]bool{},
		rooms:  map[string]components.TemperatureControllerInternal{},
		trvs:   map[string][]components.TemperatureController{},
	}
}

func (t *tempd) Name() string {
	return "tempd"
}

func (t *tempd) Init(config *config.Config) error {
	t.config = config.TemperatureControl
	return nil
}

func (t *tempd) addTrv(roomName string, tc components.TemperatureController) {
	if len(t.trvs[roomName]) == 0 {
		t.trvs[roomName] = []components.TemperatureController{}
	}

	t.trvs[roomName] = append(t.trvs[roomName], tc)
}

func (t *tempd) init() {
	for _, component := range t.components.List() {
		roomName := component.Device().Room
		switch component.Type() {
		case components.TypeBoiler:
			t.boiler = component.(components.Switch)
		case components.TypeHomedTemperature:
			t.rooms[roomName] = component.(components.TemperatureControllerInternal)
		case components.TypeZigbeeTRV:
			t.addTrv(roomName, component.(components.TemperatureController))
		}
	}
}

func (t *tempd) run() {
	t.updateTemperatureValues()
	t.updateTemperatureMode()
	t.setTRVTarget()
	t.setBoilerState()
}

func (t *tempd) Run(ctx context.Context, config *apps.Config) error {
	t.logger = config.Logger.With(zap.String("app", t.Name()))

	if !t.config.Enabled {
		t.logger.Info("app is disabled")
		return nil
	}

	t.components = config.Components
	t.init()

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			t.run()
		}
	}
}

// Temperature returns the temperature in the room
func (t *tempd) roomTemperature(room string) float64 {
	var temperature float64
	var found float64

	cs := t.components.ListByRoom(room)
	for _, c := range cs {
		if c.Type() == components.TypeHomedTemperature {
			continue
		}

		if c.Type() == components.TypeZigbeeTRV {
			// Don't use this for now
			continue
		}

		tc, ok := c.(components.TemperatureGetter)
		if !ok {
			continue
		}

		logger := c.LoggerWithFields(t.logger)

		v, err := tc.Temperature()
		if err != nil {
			logger.Warn("failed to get room temperature from component", zap.Error(err))
			continue
		}
		found++
		temperature += v
	}

	if found == 0 {
		return 0
	}

	return (temperature / found)
}

// This function get the temperature from multiple devices in the room and set
// the homed's own temperature
func (t *tempd) updateTemperatureValues() {
	for room, component := range t.rooms {
		_, ok := t.rooms[room]
		if !ok {
			t.logger.Info("failed to find room", zap.String("room", room))
			continue
		}

		// Update the current room temperature
		currentTemperature := t.roomTemperature(room)
		if math.IsNaN(currentTemperature) {
			t.logger.Warn("temperature is NaN")
			continue
		}

		t.logger.Debug(
			"setting room temperature",
			zap.String("room", room),
			zap.Float64("temperature", currentTemperature),
		)
		if err := component.SetTemperature(currentTemperature); err != nil {
			t.logger.Warn(
				"failed to update room temperature",
				zap.Error(err),
				zap.String("room", room),
			)
			continue
		}
	}
}

func (t *tempd) setBoilerState() {
	if t.boiler == nil {
		return
	}

	boiler := t.boiler

	shouldTurnOn := false

	for room, controller := range t.rooms {
		logger := controller.LoggerWithFields(t.logger)

		_, ok := t.rising[room]
		if !ok {
			t.rising[room] = false
		}

		current, err := controller.Temperature()
		if err != nil {
			logger.Warn("failed to get temperature", zap.Error(err))
			continue
		}

		if current == 0 {
			logger.Info("the temperature is reported to be 0, this is unlikely, let's ignore this for now")
			continue
		}

		target, err := controller.TemperatureTarget()
		if err != nil {
			logger.Warn("failed to get temperature target", zap.Error(err))
			continue
		}

		if target == 0 {
			logger.Info("the temperature target is 0, this is unlikely, let's ignore this for now")
			continue
		}

		min := target - t.config.Hysteresis
		max := target + t.config.Hysteresis

		logger = logger.With(
			zap.Float64("temperature", current),
			zap.Float64("target", target),
		)

		if current > max {
			if t.rising[room] {
				logger.Info("boiler entering the falling phase")
			}
			t.rising[room] = false
		}

		if current < min {
			if !t.rising[room] {
				logger.Info("boiler entering the rising phase")
			}
			t.rising[room] = true
		}

		if t.rising[room] && (current < max) {
			logger.Debug("boiler should be on")
			shouldTurnOn = true
			break
		}
	}

	logger := t.logger.With(zap.Bool("state", shouldTurnOn))

	if boiler.IsOn() == shouldTurnOn {
		logger.Debug("boiler already in the good state")
		return
	}

	logger.Info("changing boiler state")

	var err error
	if shouldTurnOn {
		err = t.boiler.TurnOn()
	} else {
		err = t.boiler.TurnOff()
	}
	if err != nil {
		logger.Warn("failed to set boiler state", zap.Error(err))
	}
}

func (t *tempd) recalibrateTRV(trv components.TemperatureController,
	roomTemperature, trvTemperature float64) error {

	logger := trv.LoggerWithFields(t.logger)

	calibration, err := trv.TemperatureCalibration()
	if err != nil {
		if errors.Is(err, components.ErrOperatingInProgress) {
			logger.Debug("calibration already in progress")
			return nil
		}

		return err
	}

	trvMesuredTemperature := trvTemperature - calibration

	delta := roomTemperature - trvMesuredTemperature

	// Only keep on decimal of precision
	delta = math.Round(delta*10) / 10

	if math.Abs(delta) > t.config.CalibrationMaxOffset {
		logger.Debug("invalid calibration, ignoring for now", zap.Float64("calibration", delta))
		return nil
	}

	if math.Abs(calibration-delta) > t.config.CalibrationThreshold {
		logger.Info("recalibrating the trv",
			zap.Float64("old_calibration", calibration),
			zap.Float64("new_calibration", delta),
			zap.Float64("calibration_diff", math.Abs(calibration-delta)),
			zap.Float64("calibration_threshold",
				t.config.CalibrationThreshold),
		)
		err = trv.SetTemperatureCalibration(delta)
		if err != nil {
			if errors.Is(err, components.ErrOperatingInProgress) {
				logger.Debug("operation already in progress")
			} else {
				logger.Warn("failed to calibrate trv", zap.Error(err))
			}
		}
	}

	return nil
}

func (t *tempd) setTRVTarget() {
	for room, controller := range t.rooms {
		logger := t.logger.With(zap.String("room", room))

		target, err := controller.TemperatureTarget()
		if err != nil {
			logger.Warn("failed to get temperature target", zap.Error(err))
			continue
		}

		roomTemperature, err := controller.Temperature()
		if err != nil {
			logger.Warn("failed to get temperature", zap.Error(err))
			continue
		}

		trvs, ok := t.trvs[room]
		if !ok {
			logger.Warn("failed to get trvs", zap.Error(err))
			continue
		}

		for _, trv := range trvs {
			trvLogger := trv.LoggerWithFields(logger)
			// Compute the new calibration
			trvTemperature, err := trv.Temperature()
			if err != nil {
				trvLogger.Warn("failed to get trv temperature", zap.Error(err))
				continue
			}

			if t.config.CalibrateTRV {
				if err := t.recalibrateTRV(trv, roomTemperature, trvTemperature); err != nil {
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
}

func (t *tempd) updateTemperatureMode() {
	for _, component := range t.rooms {
		schedule := component.Schedule()
		scheduledTarget := schedule.Value()

		logger := component.LoggerWithFields(t.logger)

		currentTarget, err := component.TemperatureTarget()
		if err != nil {
			logger.Warn("failed to get the current temperature target", zap.Error(err))
		}

		if currentTarget != scheduledTarget {
			err := component.SetTemperatureTarget(scheduledTarget)
			if err != nil {
				logger.Warn("failed to set the temperature target", zap.Error(err))
				continue
			}
		}

		mode, err := component.TemperatureMode()
		if err != nil {
			t.logger.Warn("failed to get the temperature mode", zap.Error(err))
			continue
		}

		if mode == components.TemperatureModeUntilNextChange {
			nextChange := schedule.NextChange()
			if nextChange != nil {
				if err := component.SetTemperatureModeManualUntil(nextChange); err != nil {
					logger.Warn("failed to set the temperature manual mode until", zap.Error(err))
					continue
				}
			}
		}

		manualUntil, err := component.TemperatureModeManualUntil()
		if err != nil {
			logger.Warn("failed to get the end date of the manual mode", zap.Error(err))
			continue
		}

		if manualUntil != nil && time.Now().After(*manualUntil) {
			if err := component.SetTemperatureMode(components.TemperatureModeAuto); err != nil {
				logger.Warn("failed to set temperature mode", zap.Error(err))
				continue
			}

			if err := component.SetTemperatureModeManualUntil(nil); err != nil {
				logger.Warn("failed to reset the override mode", zap.Error(err))
				continue
			}

			previousTarget, err := component.TemperatureTarget()
			if err != nil {
				logger.Warn("failed to get the temperature target", zap.Error(err))
				continue
			}

			if err := component.SetTemperatureTarget(previousTarget); err != nil {
				logger.Warn("failed to set temperature target", zap.Error(err))
				continue
			}
		}

		if err := component.PublishState(); err != nil {
			t.logger.Warn("failed to publish state", zap.Error(err))
			continue
		}
	}
}
