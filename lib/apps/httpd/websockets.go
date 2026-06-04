package httpd

import (
	"log/slog"
	"maps"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/gregdel/homed/lib/components"
)

type websocketClient struct {
	remote string
	mu     sync.Mutex
}

func (h *httpd) registerWebsocket(ws *websocket.Conn, remote string) *websocketClient {
	client := &websocketClient{remote: remote}

	h.mu.Lock()
	h.websockets[ws] = client
	h.mu.Unlock()

	h.logger.Info("new websocket registered", slog.String("remote", remote))
	return client
}

func (h *httpd) unregisterWebsocket(ws *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	client, ok := h.websockets[ws]
	if !ok {
		return
	}

	_ = ws.Close()
	delete(h.websockets, ws)
	h.logger.Info("websocket unregistered", slog.String("remote", client.remote))
}

func (h *httpd) publishToWebsocket(id string) {
	component, err := h.components.Get(id)
	if err != nil {
		h.logger.Error("failed to get component",
			slog.String("id", id))
		return
	}

	conns := map[*websocket.Conn]*websocketClient{}
	h.mu.RLock()
	maps.Copy(conns, h.websockets)
	h.mu.RUnlock()

	var wg sync.WaitGroup
	data := components.NewComponentJSON(component, h.components.HasGraph(id))
	for ws, client := range conns {
		wg.Add(1)
		go func(ws *websocket.Conn, client *websocketClient) {
			defer wg.Done()

			client.mu.Lock()
			defer client.mu.Unlock()

			if err := ws.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				h.logger.Info(
					"failed to set websocket write deadline",
					slog.String("remote", client.remote),
					slog.String("event_id", id),
					slog.Any("error", err),
				)
				h.unregisterWebsocket(ws)
				return
			}
			err := ws.WriteJSON(data)
			if err != nil {
				h.logger.Info(
					"failed to publish to websocket",
					slog.String("remote", client.remote),
					slog.String("event_id", id),
					slog.Any("error", err),
				)
				h.unregisterWebsocket(ws)
				return
			}
			if err := ws.SetWriteDeadline(time.Time{}); err != nil {
				h.logger.Info(
					"failed to clear websocket write deadline",
					slog.String("remote", client.remote),
					slog.String("event_id", id),
					slog.Any("error", err),
				)
				h.unregisterWebsocket(ws)
			}
		}(ws, client)
	}

	wg.Wait()
}

func (h *httpd) pingWebsocket(ws *websocket.Conn, client *websocketClient) error {
	client.mu.Lock()
	defer client.mu.Unlock()

	if err := ws.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		return err
	}
	if err := ws.WriteMessage(websocket.PingMessage, nil); err != nil {
		return err
	}
	return ws.SetWriteDeadline(time.Time{})
}
