package tempd

import (
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
	tuyas map[string][]components.TemperatureController
}

func app() *tempd {
	return &tempd{
		rising: map[string]bool{},
		rooms:  map[string]components.TemperatureControllerInternal{},
		tuyas:  map[string][]components.TemperatureController{},
	}
}

func (t *tempd) Name() string {
	return name
}

func (t *tempd) Init(config *config.Config) error {
	t.enabled = config.TemperatureControl
	return nil
}

func (t *tempd) init() {
	for _, component := range t.components.List() {
		roomName := component.Room()
		switch component.Type() {
		case components.TypeBoiler:
			t.boiler = component.(components.Switch)
		case components.TypeHomedTemperature:
			t.rooms[roomName] = component.(components.TemperatureControllerInternal)
		case components.TypeTuyaTRV:
			if len(t.tuyas[roomName]) == 0 {
				t.tuyas[roomName] = []components.TemperatureController{}
			}
			t.tuyas[roomName] = append(
				t.tuyas[roomName],
				component.(components.TemperatureController),
			)
		}
	}
}

func (t *tempd) Run(ctx *apps.RunCtx) error {
	if !t.enabled {
		ctx.Logger.Info("app is disabled", zap.String("app_name", name))
		return nil
	}

	t.logger = ctx.Logger
	t.components = ctx.Components
	t.init()

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	// TODO: remove this first round
	t.updateRoomsTemperatures()
	t.setRoomsTemperatures()
	t.setRoomsDevicesTemperatures()
	t.setBoilerState()

	exit := false
	for {
		if exit {
			break
		}

		select {
		case <-ctx.Ctx.Done():
			exit = true
			break
		case <-ticker.C:
			t.updateRoomsTemperatures()
			t.setRoomsTemperatures()
			t.setRoomsDevicesTemperatures()
			t.setBoilerState()
		}
	}

	return nil
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

		if c.Type() == components.TypeTuyaTRV {
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
			t.logger.Info(
				"failed to get temperature",
				zap.String("component_type", string(component.Type())),
			)
			continue
		}

		target, err := component.TemperatureTarget()
		if err != nil {
			t.logger.Info(
				"failed to get temperature target",
				zap.String("component_type", string(component.Type())),
				zap.String("component_id", component.ID()),
			)
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

		tuyas, ok := t.tuyas[roomName]
		if !ok {
			continue
		}

		current, _ := component.Temperature()
		for _, tuya := range tuyas {
			// Compute the new calibration
			computedTemp, _ := tuya.Temperature()
			calibration, _ := tuya.TemperatureCalibration()
			temp := computedTemp - calibration
			// Only keep on decimal of precision
			newCalibration := current - temp
			allowedError := 0.3

			if newCalibration < (calibration-allowedError) || newCalibration > (calibration+allowedError) {
				t.logger.Info(
					"recalibrating the tuya",
					zap.String("room", roomName),
					zap.Float64("calibration", newCalibration),
				)
				err = tuya.SetTemperatureCalibration(newCalibration)
				if err != nil {
					t.logger.Error(
						"failed to calibrate tuya",
						zap.Error(err),
					)
				}

			}

			tuyaTarget, err := tuya.TemperatureTarget()
			if err != nil {
				t.logger.Error(
					"failed to get tuya's temperature target",
					zap.Error(err),
				)
				continue
			}

			if tuyaTarget == target {
				continue
			}

			t.logger.Info(
				"should configure the tuya to the target",
				zap.String("room", roomName),
				zap.Float64("target", target),
			)

			err = tuya.SetTemperatureTarget(target)
			if err != nil {
				t.logger.Error(
					"failed to set tuya's temperature target",
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
				zap.String("device", component.Device()),
			)
			continue
		}

		mode, err := component.TemperatureMode()
		if err != nil {
			t.logger.Error(
				"failed to get the temperature mode",
				zap.Error(err),
				zap.String("room", roomName),
				zap.String("device", component.Device()),
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
						zap.String("device", component.Device()),
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
				zap.String("device", component.Device()),
			)
			continue
		}

		if manualUntil != nil && time.Now().After(*manualUntil) {
			if err := component.SetTemperatureMode(components.TemperatureModeAuto); err != nil {
				t.logger.Error(
					"failed to set temperature mode",
					zap.Error(err),
					zap.String("room", roomName),
					zap.String("device", component.Device()),
				)
				continue
			}

			if err := component.SetTemperatureModeManualUntil(nil); err != nil {
				t.logger.Error(
					"failed to set temperature target",
					zap.Error(err),
					zap.String("room", roomName),
					zap.String("device", component.Device()),
				)
				continue
			}

			previousTarget, err := component.TemperatureTarget()
			t.logger.Error(
				"failed to get the temperature target",
				zap.Error(err),
				zap.String("room", roomName),
				zap.String("device", component.Device()),
			)

			if err := component.SetTemperatureTarget(previousTarget); err != nil {
				t.logger.Error(
					"failed to set temperature target",
					zap.Error(err),
					zap.String("room", roomName),
					zap.String("device", component.Device()),
				)
				continue
			}
		}

		if err := component.PublishState(); err != nil {
			t.logger.Error(
				"failed to publish state",
				zap.Error(err),
				zap.String("room", roomName),
				zap.String("device", component.Device()),
			)
			continue
		}
	}
}
