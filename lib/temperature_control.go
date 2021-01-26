package homed

// Update the computed temperature of every room
// Check and update the target temperature from a schedule
// Check if the boiler state needs to be up

// If the target is manually updated:
// manual_until: (compute the next scheduled change || manual)
type temperatureController struct {
	// boiler *components.Boiler

	// // Wether we're in a rising or falling phase of the hysteresis algorithm
	// rising map[string]bool

	// rooms     map[string]*components.HomedTemperature
	// schedules map[string]*schedule.Schedule
	// tuyas     map[string][]*components.TuyaTRV
}

func newTemperatureController() *temperatureController {
	return &temperatureController{
		// rising:    map[string]bool{},
		// rooms:     map[string]*components.HomedTemperature{},
		// tuyas:     map[string][]*components.TuyaTRV{},
		// schedules: map[string]*schedule.Schedule{},
	}
}

// func (h *Homed) initTemperatureController() error {
// 	tc := newTemperatureController()
// 	for _, room := range h.rooms {
// 		for _, device := range room.Devices {
// 			for _, component := range device.Components {
// 				switch component.Type() {
// 				case components.TypeBoiler:
// 					tc.boiler = component.(*components.Boiler)
// 				case components.TypeHomedTemperature:
// 					tc.schedules[room.Name] = schedule.New()
// 					tc.rooms[room.Name] = component.(*components.HomedTemperature)
// 				case components.TypeTuyaTRV:
// 					if len(tc.tuyas[room.Name]) == 0 {
// 						tc.tuyas[room.Name] = []*components.TuyaTRV{}
// 					}
// 					tc.tuyas[room.Name] = append(
// 						tc.tuyas[room.Name],
// 						component.(*components.TuyaTRV),
// 					)
// 				}
// 			}
// 		}
// 	}

// 	h.temperatureController = tc

// 	h.loadTemperatureSchedules()

// 	return nil
// }

// func (h *Homed) startTemperatureControl(done <-chan struct{}) {
// 	fields := zap.String("app", "temperature_control")
// 	ticker := time.NewTicker(time.Minute)
// 	defer ticker.Stop()

// 	h.logger.Info("Starting the temperature controller", fields)

// 	// TODO: remove this first round
// 	h.updateRoomsTemperatures()
// 	h.setRoomsTemperatures()
// 	h.setRoomsDevicesTemperatures()
// 	h.setBoilerState()

// 	exit := false
// 	for {
// 		if exit {
// 			break
// 		}

// 		select {
// 		case <-done:
// 			exit = true
// 			break
// 		case <-ticker.C:
// 			h.updateRoomsTemperatures()
// 			h.setRoomsTemperatures()
// 			h.setRoomsDevicesTemperatures()
// 			h.setBoilerState()
// 		}
// 	}

// 	h.logger.Info("Exiting temperature control", fields)
// }

// func (h *Homed) updateRoomsTemperatures() {
// 	for roomName, component := range h.temperatureController.rooms {
// 		room, ok := h.rooms[roomName]
// 		if !ok {
// 			h.logger.Info("failed to find room", zap.String("room", roomName))
// 			continue
// 		}

// 		// Update the current room temperature
// 		currentTemperature := room.Temperature()
// 		if currentTemperature == 0 {
// 			continue
// 		}
// 		component.Current = currentTemperature

// 		if err := component.PublishState(h.mqttClient); err != nil {
// 			h.logger.Warn(err.Error(), zap.String("room", roomName))
// 			continue
// 		}
// 	}
// }

// func (h *Homed) setBoilerState() {
// 	if h.temperatureController.boiler == nil {
// 		return
// 	}

// 	boiler := h.temperatureController.boiler

// 	expectedBoilerState := false
// 	hysteresis := 0.3

// 	for room, component := range h.temperatureController.rooms {
// 		_, ok := h.temperatureController.rising[room]
// 		if !ok {
// 			h.temperatureController.rising[room] = false
// 		}

// 		current := component.Current
// 		target := component.CurrentTarget()
// 		min := target - hysteresis
// 		max := target + hysteresis

// 		if current > max {
// 			h.logger.Info("boiler in rising phase", zap.String("room", room))
// 			h.temperatureController.rising[room] = false
// 		}

// 		if current < min {
// 			h.logger.Info("boiler in falling phase", zap.String("room", room))
// 			h.temperatureController.rising[room] = true
// 		}

// 		if h.temperatureController.rising[room] && (current < max) {
// 			h.logger.Info("boiler should be on", zap.String("room", room))
// 			expectedBoilerState = true
// 			break
// 		}
// 	}

// 	if boiler.On == expectedBoilerState {
// 		h.logger.Info("boiler already in the good state", zap.Bool("state", expectedBoilerState))
// 		return
// 	}

// 	h.logger.Info("changing boiler state", zap.Bool("state", expectedBoilerState))
// 	data := []byte("OFF")
// 	if expectedBoilerState {
// 		data = []byte("ON")
// 	}

// 	err := h.temperatureController.boiler.WriteCommand(h.mqttClient, data)
// 	if err != nil {
// 		h.logger.Error(err.Error())
// 	}
// }

// func (h *Homed) setRoomsDevicesTemperatures() {
// 	for roomName, component := range h.temperatureController.rooms {
// 		target := component.CurrentTarget()

// 		tuyas, ok := h.temperatureController.tuyas[roomName]
// 		if !ok {
// 			continue
// 		}

// 		// data, err := json.Marshal(components.TuyaTRV{HeatingSetpoint: component.Target})
// 		// if err != nil {
// 		// 	h.logger.Warn(err.Error())
// 		// 	continue
// 		// }

// 		for _, tuya := range tuyas {
// 			if tuya.HeatingSetpoint == target {
// 				continue
// 			}

// 			h.logger.Info(
// 				"should configure the tuya to the target",
// 				zap.String("room", roomName),
// 				zap.Float64("target", target),
// 			)

// 			// err := tuya.WriteCommand(h.mqttClient, data)
// 			// if err != nil {
// 			// 	h.logger.Warn(err.Error())
// 			// 	continue
// 			// }
// 		}
// 	}
// }

// func (h *Homed) setRoomsTemperatures() {
// 	for roomName, component := range h.temperatureController.rooms {
// 		// Update the target state
// 		component.Target = h.temperatureTarget(roomName)

// 		if component.Mode == components.HomedTemperatureModeNextTimeBlock {
// 			nextTime := h.timeOfNextTimeBlock(roomName)
// 			if nextTime != nil {
// 				component.ManualUntil = nextTime
// 			}
// 		}

// 		if component.ManualUntil != nil && time.Now().After(*component.ManualUntil) {
// 			component.Mode = components.HomedTemperatureModeAuto
// 			component.ManualUntil = nil
// 			component.ManualTarget = component.Target
// 		}

// 		if err := component.PublishState(h.mqttClient); err != nil {
// 			h.logger.Warn(err.Error(), zap.String("room", roomName))
// 			continue
// 		}
// 	}
// }

// func (h *Homed) timeOfNextTimeBlock(room string) *time.Time {
// 	schedule, ok := h.temperatureController.schedules[room]
// 	if !ok {
// 		h.logger.Error("missing schedule for room", zap.String("room", room))
// 		return nil
// 	}

// 	_, t := schedule.NextTime()
// 	return t
// }

// func (h *Homed) temperatureTarget(room string) float64 {
// 	var defaultTarget float64 = 14
// 	schedule, ok := h.temperatureController.schedules[room]
// 	if !ok {
// 		h.logger.Error("missing schedule for room", zap.String("room", room))
// 		return defaultTarget
// 	}

// 	ts := schedule.Now()
// 	if ts == nil {
// 		return defaultTarget
// 	}

// 	return ts.Value
// }

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
