package httpd

import (
	"github.com/gorilla/websocket"
	"github.com/gregdel/homed/lib/components"
	"go.uber.org/zap"
)

func (h *httpd) registerWebsocket(ws *websocket.Conn, remote string) {
	h.mu.Lock()
	h.websockets[ws] = remote
	h.mu.Unlock()

	h.logger.Info("new websocket registered", zap.String("remote", remote))
}

func (h *httpd) unregisterWebsocket(ws *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	remote, ok := h.websockets[ws]
	if !ok {
		return
	}

	ws.Close()
	delete(h.websockets, ws)
	h.logger.Info("websocket unregistered", zap.String("remote", remote))
}

func (h *httpd) publishToWebsocket(id string) {
	component, err := h.components.Get(id)
	if err != nil {
		h.logger.Error("failed to get component",
			zap.String("id", id))
		return
	}

	conns := map[*websocket.Conn]string{}
	h.mu.RLock()
	for ws, remote := range h.websockets {
		conns[ws] = remote
	}
	h.mu.RUnlock()

	data := components.NewComponentJSON(component)
	for ws, remote := range conns {
		err := ws.WriteJSON(data)
		if err != nil {
			h.logger.Info(
				"failed to publish to websocket",
				zap.String("remote", remote),
				zap.String("event_id", id),
				zap.Error(err),
			)
			h.unregisterWebsocket(ws)
			continue
		}
	}
}
