package homed

import (
	"time"

	"github.com/gregdel/homed/lib/components"
	"go.uber.org/zap"
)

// If the target is manually updated:
// manual_until: (compute the next scheduled change || manual)
type temperatureController struct {
	boiler components.Switch

	// Wether we're in a rising or falling phase of the hysteresis algorithm
	rising map[string]bool

	rooms map[string]components.TemperatureControllerInternal
	tuyas map[string][]components.TemperatureController
}

func newTemperatureController() *temperatureController {
	return &temperatureController{
		rising: map[string]bool{},
		rooms:  map[string]components.TemperatureControllerInternal{},
		tuyas:  map[string][]components.TemperatureController{},
	}
}

func (h *Homed) initTemperatureController() error {
	tc := newTemperatureController()
	for _, room := range h.rooms {
		for _, component := range room.Components() {
			switch component.Type() {
			case components.TypeBoiler:
				tc.boiler = component.(components.Switch)
			case components.TypeHomedTemperature:
				tc.rooms[room.Name] = component.(components.TemperatureControllerInternal)
			case components.TypeTuyaTRV:
				if len(tc.tuyas[room.Name]) == 0 {
					tc.tuyas[room.Name] = []components.TemperatureController{}
				}
				tc.tuyas[room.Name] = append(
					tc.tuyas[room.Name],
					component.(components.TemperatureController),
				)
			}
		}
	}

	h.temperatureController = tc

	// h.loadTemperatureSchedules()

	return nil
}

func (h *Homed) startTemperatureControl(done <-chan struct{}) {
	fields := zap.String("app", "temperature_control")
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	h.logger.Info("Starting the temperature controller", fields)

	// TODO: remove this first round
	h.updateRoomsTemperatures()
	h.setRoomsTemperatures()
	h.setRoomsDevicesTemperatures()
	h.setBoilerState()

	exit := false
	for {
		if exit {
			break
		}

		select {
		case <-done:
			exit = true
			break
		case <-ticker.C:
			h.updateRoomsTemperatures()
			h.setRoomsTemperatures()
			h.setRoomsDevicesTemperatures()
			h.setBoilerState()
		}
	}

	h.logger.Info("Exiting temperature control", fields)
}

func (h *Homed) updateRoomsTemperatures() {
	for roomName, component := range h.temperatureController.rooms {
		room, ok := h.rooms[roomName]
		if !ok {
			h.logger.Info("failed to find room", zap.String("room", roomName))
			continue
		}

		// Update the current room temperature
		currentTemperature := room.Temperature()
		if currentTemperature == 0 {
			continue
		}

		if err := component.SetTemperature(currentTemperature); err != nil {
			h.logger.Error(err.Error(), zap.String("room", roomName))
			continue
		}
	}
}

func (h *Homed) setBoilerState() {
	if h.temperatureController.boiler == nil {
		return
	}

	boiler := h.temperatureController.boiler

	expectedBoilerState := false
	hysteresis := 0.3

	for room, component := range h.temperatureController.rooms {
		_, ok := h.temperatureController.rising[room]
		if !ok {
			h.temperatureController.rising[room] = false
		}

		current, err := component.Temperature()
		if err != nil {
			h.logger.Info(
				"failed to get temperature",
				zap.String("component_type", string(component.Type())),
			)
			continue
		}

		target, err := component.TemperatureTarget()
		if err != nil {
			h.logger.Info(
				"failed to get temperature target",
				zap.String("component_type", string(component.Type())),
				zap.String("component_id", component.ID().String()),
			)
			continue
		}

		min := target - hysteresis
		max := target + hysteresis

		if current > max {
			h.logger.Info("boiler in rising phase", zap.String("room", room))
			h.temperatureController.rising[room] = false
		}

		if current < min {
			h.logger.Info("boiler in falling phase", zap.String("room", room))
			h.temperatureController.rising[room] = true
		}

		if h.temperatureController.rising[room] && (current < max) {
			h.logger.Info("boiler should be on", zap.String("room", room))
			expectedBoilerState = true
			break
		}
	}

	if boiler.IsOn() == expectedBoilerState {
		h.logger.Info("boiler already in the good state", zap.Bool("state", expectedBoilerState))
		return
	}

	h.logger.Info("changing boiler state", zap.Bool("state", expectedBoilerState))

	err := h.temperatureController.boiler.Set(expectedBoilerState)
	if err != nil {
		h.logger.Error(err.Error())
	}
}

func (h *Homed) setRoomsDevicesTemperatures() {
	for roomName, component := range h.temperatureController.rooms {
		target, err := component.TemperatureTarget()
		if err != nil {
			h.logger.Error(
				"failed to get temperature target",
				zap.Error(err),
			)
			continue
		}

		tuyas, ok := h.temperatureController.tuyas[roomName]
		if !ok {
			continue
		}

		for _, tuya := range tuyas {
			tuyaTarget, err := tuya.TemperatureTarget()
			if err != nil {
				h.logger.Error(
					"failed to get tuya's temperature target",
					zap.Error(err),
				)
				continue
			}

			if tuyaTarget == target {
				continue
			}

			h.logger.Info(
				"should configure the tuya to the target",
				zap.String("room", roomName),
				zap.Float64("target", target),
			)

			err = tuya.SetTemperatureTarget(target)
			if err != nil {
				h.logger.Error(
					"failed to set tuya's temperature target",
					zap.Error(err),
				)
				continue
			}
		}
	}
}

func (h *Homed) setRoomsTemperatures() {
	for roomName, component := range h.temperatureController.rooms {
		scheduledTarget := component.CurrentSchedule()

		// Update the target state
		err := component.SetTemperatureTarget(scheduledTarget)
		if err != nil {
			h.logger.Error(
				"failed to set the temperature target",
				zap.Error(err),
			)
			continue
		}

		mode, err := component.TemperatureMode()
		if err != nil {
			h.logger.Error(
				"failed to get the temperature mode",
				zap.Error(err),
			)
			continue
		}

		if mode == components.TemperatureModeNextTimeBlock {
			nextTime := component.ScheduledNextTime()
			if nextTime != nil {
				if err := component.SetTemperatureModeManualUntil(nextTime); err != nil {
					h.logger.Error(
						"failed to set the temperature manual mode until",
						zap.Error(err),
					)
					continue
				}
			}
		}

		manualUntil, err := component.TemperatureModeManualUntil()
		if err != nil {
			h.logger.Error(
				"failed to get the end date of the manual mode",
				zap.Error(err),
			)
			continue
		}

		if manualUntil != nil && time.Now().After(*manualUntil) {
			if err := component.SetTemperatureMode(components.TemperatureModeAuto); err != nil {
				h.logger.Error(
					"failed to set temperature mode",
					zap.Error(err),
				)
				continue
			}

			if err := component.SetTemperatureModeManualUntil(nil); err != nil {
				h.logger.Error(
					"failed to set temperature target",
					zap.Error(err),
				)
				continue
			}

			previousTarget, err := component.TemperatureTarget()
			h.logger.Error(
				"failed to get the temperature target",
				zap.Error(err),
			)

			if err := component.SetTemperatureTarget(previousTarget); err != nil {
				h.logger.Error(
					"failed to set temperature target",
					zap.Error(err),
				)
				continue
			}
		}

		if err := component.PublishState(); err != nil {
			h.logger.Error(
				"failed to publish state",
				zap.Error(err),
				zap.String("room", roomName),
			)
			continue
		}
	}
}

// // TODO

// func (h *Homed) loadTemperatureSchedules() {
// 	err := readFile(h.scheduleFile, h.temperatureController.schedules)
// 	if err != nil {
// 		h.logger.Error(
// 			"failed to read temperature schedules",
// 			zap.String("error", err.Error()),
// 		)
// 	}
// }

// func (h *Homed) saveTemperatureSchedules() {
// 	err := writeFile(h.scheduleFile, true, h.temperatureController.schedules)
// 	if err != nil {
// 		h.logger.Error(
// 			"failed to save temperature schedules",
// 			zap.String("error", err.Error()),
// 		)
// 	}
// }
