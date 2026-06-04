package httpd

import (
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

// Time allowed to read the next pong message from the client
const pongWait = 15 * time.Second

// Time allowed to write to a websocket
const writeWait = 15 * time.Second

// Send pings before the read deadline expires.
const pingPeriod = (pongWait * 9) / 10

func (h *httpd) httpRender(w http.ResponseWriter, status string, data any) {
	o := struct {
		Status string `json:"status"`
		Data   any    `json:"data"`
	}{
		Status: status,
		Data:   data,
	}

	if err := h.render.JSON(w, http.StatusOK, o); err != nil {
		h.httpError(w, fmt.Sprintf("failed to render json data: %s", err))
	}
}

func (h *httpd) httpRenderJSON(w http.ResponseWriter, data any) {
	h.httpRender(w, "success", data)
}

func (h *httpd) httpError(w http.ResponseWriter, msg string) {
	h.httpRender(w, "error", msg)
}

func (h *httpd) httpComponentList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	h.httpRenderJSON(w, h.components)
}

func (h *httpd) updateComponent(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	component, err := h.components.Get(id)
	if err != nil {
		h.httpError(w, fmt.Sprintf("failed to get the component: %s", err))
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		h.httpError(w, fmt.Sprintf("failed to read body: %s", err))
		return
	}

	mqttClient := component.MQTTClient()
	if mqttClient != nil && !mqttClient.IsConnectionOpen() {
		h.httpError(w, "not connected to the mqtt broker")
		return
	}

	err = component.WriteCommand(data)
	if err != nil {
		h.logger.Warn("failed to write mqtt command", slog.Any("error", err))
		h.httpError(w, fmt.Sprintf("failed to write mqtt command: %s", err))
		return
	}

	h.httpRenderJSON(w, nil)
}

func (h *httpd) websocketEvents(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	host, port, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		h.httpError(w, fmt.Sprintf("failed to get remote addr: %s", err))
		return
	}

	if r.Header.Get("X-Real-IP") != "" {
		host = r.Header.Get("X-Real-IP")
	}

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	// Upgrade the request for websockets
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.httpError(w, fmt.Sprintf("got error upgrading request: %s", err))
		return
	}

	// The pong handler only postpone the read deadline
	_ = ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error {
		_ = ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	client := h.registerWebsocket(ws, net.JoinHostPort(host, port))
	done := make(chan struct{})
	defer h.unregisterWebsocket(ws)
	defer close(done)

	go func() {
		ticker := time.NewTicker(pingPeriod)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				if err := h.pingWebsocket(ws, client); err != nil {
					h.logger.Info(
						"failed to ping websocket",
						slog.String("remote", client.remote),
						slog.Any("error", err),
					)
					h.unregisterWebsocket(ws)
					return
				}
			}
		}
	}()

	for !h.isExiting() {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}
}
