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

const name = "tempd"

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
	return name
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
	if !t.config.Enabled {
		config.Logger.Info("app is disabled", zap.String("app_name", name))
		return nil
	}

	t.logger = config.Logger
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
	var temperature float64 = 0
	var found float64 = 0

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

		v, err := tc.Temperature()
		if err != nil {
			t.logger.Warn("failed to get room temperature from component",
				zap.String("room", room),
				zap.String("friendly_name", c.FriendlyName()),
				zap.Error(err))
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
			t.logger.Warn("tempd: temperature is NaN")
			continue
		}

		t.logger.Debug(
			"Setting homed's room temperature",
			zap.String("room", room),
			zap.Float64("temperature", currentTemperature),
		)
		if err := component.SetTemperature(currentTemperature); err != nil {
			t.logger.Warn(err.Error(), zap.String("room", room))
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
		_, ok := t.rising[room]
		if !ok {
			t.rising[room] = false
		}

		current, err := controller.Temperature()
		if err != nil {
			t.logger.Warn("failed to get temperature", zap.Error(err))
			continue
		}

		if current == 0 {
			t.logger.Info("the temperature is reported to be 0, this is unlikely, let's ignore this for now")
			continue
		}

		target, err := controller.TemperatureTarget()
		if err != nil {
			t.logger.Warn(
				"failed to get temperature target",
				zap.String("component_type", string(controller.Type())),
				zap.String("component_id", controller.ID()),
				zap.Error(err),
			)
			continue
		}

		if target == 0 {
			t.logger.Info("the temperature target is 0, this is unlikely, let's ignore this for now")
			continue
		}

		min := target - t.config.Hysteresis
		max := target + t.config.Hysteresis

		zapFields := []zap.Field{
			zap.String("room", room),
			zap.Float64("temperature", current),
			zap.Float64("target", target),
		}

		if current > max {
			if t.rising[room] {
				t.logger.Info("boiler entering the falling phase", zapFields...)
			}
			t.rising[room] = false
		}

		if current < min {
			if !t.rising[room] {
				t.logger.Info("boiler entering the rising phase", zapFields...)
			}
			t.rising[room] = true
		}

		if t.rising[room] && (current < max) {
			t.logger.Debug("boiler should be on", zapFields...)
			shouldTurnOn = true
			break
		}
	}

	if boiler.IsOn() == shouldTurnOn {
		t.logger.Debug("boiler already in the good state", zap.Bool("state", shouldTurnOn))
		return
	}

	t.logger.Info("changing boiler state", zap.Bool("state", shouldTurnOn))

	var err error
	if shouldTurnOn {
		err = t.boiler.TurnOn()
	} else {
		err = t.boiler.TurnOff()
	}
	if err != nil {
		t.logger.Warn("failed to set boiler state",
			zap.Bool("state", shouldTurnOn), zap.Error(err))
	}
}

func (t *tempd) recalibrateTRV(trv components.TemperatureController,
	roomTemperature, trvTemperature float64, zapFields []zap.Field) error {
	calibration, err := trv.TemperatureCalibration()
	if err != nil {
		return err
	}

	trvMesuredTemperature := trvTemperature - calibration

	delta := roomTemperature - trvMesuredTemperature

	// Only keep on decimal of precision
	delta = math.Round(delta*10) / 10

	if math.Abs(delta) > t.config.CalibrationMaxOffset {
		t.logger.Info(
			"invalid calibration, resetting calibration to 0",
			append(zapFields, zap.Float64("calibration", delta))...)
		delta = 0
	}

	if math.Abs(calibration-delta) > t.config.CalibrationThreshold {
		t.logger.Info("recalibrating the trv",
			append(zapFields,
				zap.Float64("old_calibration", calibration),
				zap.Float64("new_calibraton", delta),
				zap.Float64("calibraton_diff", math.Abs(calibration-delta)),
				zap.Float64("calibration_threshold",
					t.config.CalibrationThreshold),
			)...)
		err = trv.SetTemperatureCalibration(delta)
		if err != nil {
			if errors.Is(err, components.ErrOperatingInProgress) {
				t.logger.Debug("operation already in progress", zapFields...)
			} else {
				t.logger.Warn("failed to calibrate trv",
					append(zapFields, zap.Error(err))...)
			}
		}
	}

	return nil
}

func (t *tempd) setTRVTarget() {
	for room, controller := range t.rooms {
		target, err := controller.TemperatureTarget()
		if err != nil {
			t.logger.Warn(
				"failed to get temperature target",
				zap.Error(err),
			)
			continue
		}

		roomTemperature, err := controller.Temperature()
		if err != nil {
			t.logger.Warn("failed to get temperature", zap.Error(err))
			continue
		}

		trvs, ok := t.trvs[room]
		if !ok {
			t.logger.Warn("failed to get trvs",
				zap.String("room", room), zap.Error(err))
			continue
		}

		for _, trv := range trvs {
			zapFields := []zap.Field{
				zap.String("room", room),
				zap.String("friendly_name", trv.FriendlyName()),
			}

			// Compute the new calibration
			trvTemperature, err := trv.Temperature()
			if err != nil {
				t.logger.Warn("failed to get trv temperature",
					append(zapFields, zap.Error(err))...)
				continue
			}

			if t.config.CalibrateTRV {
				if err := t.recalibrateTRV(trv, roomTemperature, trvTemperature, zapFields); err != nil {
					t.logger.Warn("failed to calibrate the trv",
						append(zapFields, zap.Error(err))...)
					continue
				}
			}

			trvTarget, err := trv.TemperatureTarget()
			if err != nil {
				t.logger.Warn("failed to get trv's temperature target",
					append(zapFields, zap.Error(err))...)
				continue
			}

			if trvTarget == target {
				continue
			}

			t.logger.Info("should configure the trv to the target",
				append(zapFields, zap.Float64("target", target))...)

			err = trv.SetTemperatureTarget(target)
			if err != nil {
				t.logger.Warn("failed to set trv's temperature target",
					append(zapFields, zap.Error(err))...)
				continue
			}
		}
	}
}

func (t *tempd) updateTemperatureMode() {
	for roomName, component := range t.rooms {
		schedule := component.Schedule()
		scheduledTarget := schedule.Value()

		zapFields := []zap.Field{
			zap.String("room", roomName),
			zap.String("device", component.Device().Name),
		}

		currentTarget, err := component.TemperatureTarget()
		if err != nil {
			t.logger.Warn("failed to get the current temperature target",
				append(zapFields, zap.Error(err))...)
		}

		if currentTarget != scheduledTarget {
			err := component.SetTemperatureTarget(scheduledTarget)
			if err != nil {
				t.logger.Warn("failed to set the temperature target",
					append(zapFields, zap.Error(err))...)
				continue
			}
		}

		mode, err := component.TemperatureMode()
		if err != nil {
			t.logger.Warn("failed to get the temperature mode",
				append(zapFields, zap.Error(err))...)
			continue
		}

		if mode == components.TemperatureModeNextTimeBlock {
			nextTime := schedule.NextTime()
			if nextTime != nil {
				if err := component.SetTemperatureModeManualUntil(nextTime); err != nil {
					t.logger.Warn("failed to set the temperature manual mode until",
						append(zapFields, zap.Error(err))...)
					continue
				}
			}
		}

		manualUntil, err := component.TemperatureModeManualUntil()
		if err != nil {
			t.logger.Warn("failed to get the end date of the manual mode",
				append(zapFields, zap.Error(err))...)
			continue
		}

		if manualUntil != nil && time.Now().After(*manualUntil) {
			if err := component.SetTemperatureMode(components.TemperatureModeAuto); err != nil {
				t.logger.Warn("failed to set temperature mode",
					append(zapFields, zap.Error(err))...)
				continue
			}

			if err := component.SetTemperatureModeManualUntil(nil); err != nil {
				t.logger.Warn("failed to reset the override mode",
					append(zapFields, zap.Error(err))...)
				continue
			}

			previousTarget, err := component.TemperatureTarget()
			t.logger.Warn("failed to get the temperature target",
				append(zapFields, zap.Error(err))...)

			if err := component.SetTemperatureTarget(previousTarget); err != nil {
				t.logger.Warn("failed to set temperature target",
					append(zapFields, zap.Error(err))...)
				continue
			}
		}

		if err := component.PublishState(); err != nil {
			t.logger.Warn("failed to publish state",
				append(zapFields, zap.Error(err))...)
			continue
		}
	}
}
