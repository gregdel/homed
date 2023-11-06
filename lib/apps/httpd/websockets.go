package httpd

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/gregdel/homed/lib/components"
	"go.uber.org/zap"
)

func (h *httpd) registerWebsocket(ws *websocket.Conn) (string, error) {
	u, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	uuid := u.String()
	h.mu.Lock()
	h.websockets[uuid] = ws
	h.mu.Unlock()

	h.logger.Info("new websocket registered", zap.String("uuid", uuid))

	return uuid, nil
}

func (h *httpd) unregisterWebsocket(uuid string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.websockets, uuid)
	h.logger.Info("websocket unregistered", zap.String("uuid", uuid))
}

func (h *httpd) publishToWebsocket(component components.Component) {
	toUnregister := []string{}

	data := components.NewComponentJSON(component)
	h.mu.RLock()
	for uuid, ws := range h.websockets {
		err := ws.WriteJSON(data)
		if err != nil {
			h.logger.Info("failed to publish to websocket",
				zap.String("error", err.Error()),
				zap.String("uuid", uuid))
			toUnregister = append(toUnregister, uuid)
			continue
		}
	}
	h.mu.RUnlock()

	for _, uuid := range toUnregister {
		h.unregisterWebsocket(uuid)
	}
}
