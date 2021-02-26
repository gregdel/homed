package httpd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

func (h *httpd) httpRender(w http.ResponseWriter, status string, data interface{}) {
	o := struct {
		Status string      `json:"status"`
		Data   interface{} `json:"data"`
	}{
		Status: status,
		Data:   data,
	}

	if err := h.render.JSON(w, http.StatusOK, o); err != nil {
		h.httpError(w, fmt.Sprintf("failed to render json data: %s", err))
	}
}

func (h *httpd) httpRenderJSON(w http.ResponseWriter, data interface{}) {
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

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.httpError(w, fmt.Sprintf("failed to read body: %s", err))
		return
	}

	err = component.WriteCommand(data)
	if err != nil {
		h.httpError(w, fmt.Sprintf("failed to write mqtt command: %s", err))
		return
	}
}

func (h *httpd) websocketEvents(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	const (
		// Ping every 30 seconds, must be less than pongWait
		pingWait = 10 * time.Second
		// Time allowed to read the next pong message from the client
		pongWait = 15 * time.Second
		// Time allowed to write to the client
		writeWait = 10 * time.Second
	)

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
	defer ws.Close()

	// The pong handler only postpone the read deadline
	ws.SetPongHandler(func(string) error {
		_ = ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	uuid, err := h.registerWebsocket(ws)
	if err != nil {
		h.httpError(w, fmt.Sprintf("failed to register websocket: %s", err))
		return
	}
	defer h.unregisterWebsocket(uuid)

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}
}
