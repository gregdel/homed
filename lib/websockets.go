package homed

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/gregdel/homed/lib/components"
	"go.uber.org/zap"
)

func (h *Homed) registerWebsocket(ws *websocket.Conn) (string, error) {
	u, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	uuid := u.String()
	h.websockets[uuid] = ws

	h.logger.Info("new websocket registered", zap.String("uuid", uuid))

	return uuid, nil
}

func (h *Homed) unregisterWebsocket(uuid string) {
	delete(h.websockets, uuid)
	h.logger.Info("websocket unregistered", zap.String("uuid", uuid))
}

func (h *Homed) publishToWebsocket(component components.Component) {
	data := components.NewComponentJSON(component)
	for uuid, ws := range h.websockets {
		err := ws.WriteJSON(data)
		if err != nil {
			h.logger.Info("failed to publish to websocket",
				zap.String("error", err.Error()),
				zap.String("uuid", uuid))
			h.unregisterWebsocket(uuid)
			continue
		}
	}
}
