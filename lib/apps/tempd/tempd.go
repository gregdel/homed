package tempd

import (
	"context"
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

// If the target is manually updated:
// manual_until: (compute the next scheduled change || manual)
type tempd struct {
	enabled bool

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
	t.enabled = config.TemperatureControl
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
		case components.TypeSaswellTRV:
			t.addTrv(roomName, component.(components.TemperatureController))
		case components.TypeTuyaTRV:
			t.addTrv(roomName, component.(components.TemperatureController))
		}
	}
}

func (t *tempd) Run(ctx context.Context, config *apps.Config) error {
	if !t.enabled {
		config.Logger.Info("app is disabled", zap.String("app_name", name))
		return nil
	}

	t.logger = config.Logger
	t.components = config.Components
	t.init()

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	// TODO: remove this first round
	t.updateRoomsTemperatures()
	t.setRoomsTemperatures()
	t.setRoomsDevicesTemperatures()
	t.setBoilerState()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			t.updateRoomsTemperatures()
			t.setRoomsTemperatures()
			t.setRoomsDevicesTemperatures()
			t.setBoilerState()
		}
	}
}

// Temperature returns the temperature in the room
func (t *tempd) roomTemperature(room string) float64 {
	var temperature float64
	var found float64 = 0

	cs := t.components.ListByRoom(room)
	for _, c := range cs {
		if c.Type() == components.TypeHomedTemperature {
			continue
		}

		if c.Type() == components.TypeTuyaTRV || c.Type() == components.TypeSaswellTRV {
			// Don't use this for now
			continue
		}

		tc, ok := c.(components.TemperatureGetter)
		if !ok {
			continue
		}

		found = found + 1
		t, _ := tc.Temperature()
		temperature = (temperature + t) / found
	}

	return temperature
}

func (t *tempd) updateRoomsTemperatures() {
	for roomName, component := range t.rooms {
		_, ok := t.rooms[roomName]
		if !ok {
			t.logger.Info("failed to find room", zap.String("room", roomName))
			continue
		}

		// Update the current room temperature
		currentTemperature := t.roomTemperature(roomName)
		if math.IsNaN(currentTemperature) {
			t.logger.Warn("tempd: temperature is NaN")
			continue
		}

		t.logger.Debug(
			"Setting homed's room temperature",
			zap.String("room", roomName),
			zap.Float64("temperature", currentTemperature),
		)
		if err := component.SetTemperature(currentTemperature); err != nil {
			t.logger.Error(err.Error(), zap.String("room", roomName))
			continue
		}
	}
}

func (t *tempd) setBoilerState() {
	if t.boiler == nil {
		return
	}

	boiler := t.boiler

	expectedBoilerState := false
	hysteresis := 0.3

	for room, component := range t.rooms {
		_, ok := t.rising[room]
		if !ok {
			t.rising[room] = false
		}

		current, err := component.Temperature()
		if err != nil {
			t.logger.Warn(
				"failed to get temperature",
				zap.String("component_type", string(component.Type())),
				zap.Error(err),
			)
			continue
		}

		if current == 0 {
			t.logger.Info("the temperature is reported to be 0, this is unlikely, let's ignore this for now")
			continue
		}

		target, err := component.TemperatureTarget()
		if err != nil {
			t.logger.Warn(
				"failed to get temperature target",
				zap.String("component_type", string(component.Type())),
				zap.String("component_id", component.ID()),
				zap.Error(err),
			)
			continue
		}

		if target == 0 {
			t.logger.Info("the temperature target is 0, this is unlikely, let's ignore this for now")
			continue
		}

		min := target - hysteresis
		max := target + hysteresis

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
			expectedBoilerState = true
			break
		}
	}

	if boiler.IsOn() == expectedBoilerState {
		t.logger.Debug("boiler already in the good state", zap.Bool("state", expectedBoilerState))
		return
	}

	t.logger.Info("changing boiler state", zap.Bool("state", expectedBoilerState))

	err := t.boiler.Set(expectedBoilerState)
	if err != nil {
		t.logger.Error(err.Error())
	}
}

func (t *tempd) setRoomsDevicesTemperatures() {
	for roomName, component := range t.rooms {
		target, err := component.TemperatureTarget()
		if err != nil {
			t.logger.Error(
				"failed to get temperature target",
				zap.Error(err),
			)
			continue
		}

		trvs, ok := t.trvs[roomName]
		if !ok {
			continue
		}

		current, _ := component.Temperature()
		for _, trv := range trvs {
			// Compute the new calibration
			computedTemp, _ := trv.Temperature()
			calibration, _ := trv.TemperatureCalibration()
			temp := computedTemp - calibration

			newCalibration := current - temp
			allowedError := 0.3

			// Only keep on decimal of precision
			newCalibration = math.Round(newCalibration*10) / 10

			if newCalibration < (calibration-allowedError) || newCalibration > (calibration+allowedError) {
				t.logger.Info(
					"recalibrating the trv",
					zap.String("room", roomName),
					zap.Float64("calibration", newCalibration),
				)
				err = trv.SetTemperatureCalibration(newCalibration)
				if err != nil {
					t.logger.Error(
						"failed to calibrate trv",
						zap.Error(err),
					)
				}

			}

			trvTarget, err := trv.TemperatureTarget()
			if err != nil {
				t.logger.Error(
					"failed to get trv's temperature target",
					zap.Error(err),
				)
				continue
			}

			if trvTarget == target {
				continue
			}

			t.logger.Info(
				"should configure the trv to the target",
				zap.String("room", roomName),
				zap.Float64("target", target),
			)

			err = trv.SetTemperatureTarget(target)
			if err != nil {
				t.logger.Error(
					"failed to set trv's temperature target",
					zap.Error(err),
				)
				continue
			}
		}
	}
}

func (t *tempd) setRoomsTemperatures() {
	for roomName, component := range t.rooms {
		schedule := component.Schedule()
		scheduledTarget := schedule.Value()

		// Update the target state
		err := component.SetTemperatureTarget(scheduledTarget)
		if err != nil {
			t.logger.Error(
				"failed to set the temperature target",
				zap.Error(err),
				zap.String("room", roomName),
				zap.String("device", component.Device().Name),
			)
			continue
		}

		mode, err := component.TemperatureMode()
		if err != nil {
			t.logger.Error(
				"failed to get the temperature mode",
				zap.Error(err),
				zap.String("room", roomName),
				zap.String("device", component.Device().Name),
			)
			continue
		}

		if mode == components.TemperatureModeNextTimeBlock {
			nextTime := schedule.NextTime()
			if nextTime != nil {
				if err := component.SetTemperatureModeManualUntil(nextTime); err != nil {
					t.logger.Error(
						"failed to set the temperature manual mode until",
						zap.Error(err),
						zap.String("room", roomName),
						zap.String("device", component.Device().Name),
					)
					continue
				}
			}
		}

		manualUntil, err := component.TemperatureModeManualUntil()
		if err != nil {
			t.logger.Error(
				"failed to get the end date of the manual mode",
				zap.Error(err),
				zap.String("room", roomName),
				zap.String("device", component.Device().Name),
			)
			continue
		}

		if manualUntil != nil && time.Now().After(*manualUntil) {
			if err := component.SetTemperatureMode(components.TemperatureModeAuto); err != nil {
				t.logger.Error(
					"failed to set temperature mode",
					zap.Error(err),
					zap.String("room", roomName),
					zap.String("device", component.Device().Name),
				)
				continue
			}

			if err := component.SetTemperatureModeManualUntil(nil); err != nil {
				t.logger.Error(
					"failed to set temperature target",
					zap.Error(err),
					zap.String("room", roomName),
					zap.String("device", component.Device().Name),
				)
				continue
			}

			previousTarget, err := component.TemperatureTarget()
			t.logger.Error(
				"failed to get the temperature target",
				zap.Error(err),
				zap.String("room", roomName),
				zap.String("device", component.Device().Name),
			)

			if err := component.SetTemperatureTarget(previousTarget); err != nil {
				t.logger.Error(
					"failed to set temperature target",
					zap.Error(err),
					zap.String("room", roomName),
					zap.String("device", component.Device().Name),
				)
				continue
			}
		}

		if err := component.PublishState(); err != nil {
			t.logger.Error(
				"failed to publish state",
				zap.Error(err),
				zap.String("room", roomName),
				zap.String("device", component.Device().Name),
			)
			continue
		}
	}
}
