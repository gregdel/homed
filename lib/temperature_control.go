package homed

import (
	"fmt"
	"time"

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
		topic := "homed/temperature/rooms/" + room.Name + "/current"

		temperature := room.Temperature()
		if temperature == 0 {
			continue
		}
		// Publish the temperature of the room
		token := h.mqttClient.Publish(topic, 0, false, fmt.Sprintf("%f", temperature))
		if token.Wait() && token.Error() != nil {
			h.logger.Warn(token.Error().Error(), fields)
			continue
		}
	}
}
