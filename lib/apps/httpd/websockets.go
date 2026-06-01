package httpd

import (
	"log/slog"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/gregdel/homed/lib/components"
)

func (h *httpd) registerWebsocket(ws *websocket.Conn, remote string) {
	h.mu.Lock()
	h.websockets[ws] = remote
	h.mu.Unlock()

	h.logger.Info("new websocket registered", slog.String("remote", remote))
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
	h.logger.Info("websocket unregistered", slog.String("remote", remote))
}

func (h *httpd) publishToWebsocket(id string) {
	component, err := h.components.Get(id)
	if err != nil {
		h.logger.Error("failed to get component",
			slog.String("id", id))
		return
	}

	conns := map[*websocket.Conn]string{}
	h.mu.RLock()
	for ws, remote := range h.websockets {
		conns[ws] = remote
	}
	h.mu.RUnlock()

	var wg sync.WaitGroup
	data := components.NewComponentJSON(component, h.components.HasGraph(id))
	for ws, remote := range conns {
		wg.Add(1)
		go func(ws *websocket.Conn, remote string) {
			defer wg.Done()

			ws.SetWriteDeadline(time.Now().Add(writeWait))
			err := ws.WriteJSON(data)
			if err != nil {
				h.logger.Info(
					"failed to publish to websocket",
					slog.String("remote", remote),
					slog.String("event_id", id),
					slog.Any("error", err),
				)
				h.unregisterWebsocket(ws)
			}
			ws.SetWriteDeadline(time.Time{})
		}(ws, remote)
	}

	wg.Wait()
}
