package homed

import (
	"encoding/json"
	"time"

	"github.com/gregdel/homed/lib/components"
	"go.uber.org/zap"
)

func (h *Homed) temperatureControl(done <-chan struct{}) {
	fields := zap.String("app", "temperature_control")
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	h.logger.Info("Starting temperature control", fields)

	// TODO: remove this first round
	h.temperatureMQTTUpdate(fields)

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
			h.temperatureMQTTUpdate(fields)
		}
	}

	h.logger.Info("Exiting temperature control", fields)
}

func (h *Homed) temperatureMQTTUpdate(fields zap.Field) {
	for _, room := range h.rooms {
		topic := "home/homed/temperature/rooms/" + room.Name + "/state"

		temperature := room.Temperature()
		if temperature == 0 {
			continue
		}

		data := components.HomedTemperatureData{
			Current: temperature,
			Target:  h.temperatureTarget(),
		}

		payload, err := json.Marshal(data)
		if err != nil {
			h.logger.Warn(err.Error(), fields)
			continue
		}

		// Publish the temperature of the room
		token := h.mqttClient.Publish(topic, 0, true, payload)
		if token.Wait() && token.Error() != nil {
			h.logger.Warn(token.Error().Error(), fields)
			continue
		}
	}
}

func (h *Homed) temperatureTarget() float64 {
	var target float64 = 14
	var weekend bool

	now := time.Now()

	day := now.Weekday()
	if day == time.Saturday || day == time.Sunday {
		weekend = true
	}

	hour := now.Hour()
	if !weekend {
		if hour >= 10 || hour < 22 {
			return 18
		}

	} else {
		if hour >= 8 || hour < 10 {
			return 18
		}

		if hour >= 12 || hour < 13 {
			return 18
		}

		if hour >= 18 || hour < 22 {
			return 18
		}
	}

	return target
}
