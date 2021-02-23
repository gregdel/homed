package homed

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (h *Homed) initHTTP(addr string) error {
	router := httprouter.New()
	router.Handler("GET", "/metrics", promhttp.Handler())
	router.GET("/events", h.websocketEvents)

	router.GET("/components", h.httpComponentList)
	router.PUT("/components/:id", h.updateComponent)
	router.GET("/components/:id/schedule", h.httpGetSchedule)
	router.POST("/components/:id/schedule/default", h.httpPostScheduleDefault)
	router.POST("/components/:id/schedule/daily/:weekday", h.httpPostSchedule)
	router.DELETE("/components/:id/schedule/daily/:weekday/:uuid", h.httpDeleteSchedule)

	var httpFS http.FileSystem
	if h.dev {
		httpFS = http.Dir("frontend/build")
	} else {
		f, err := fs.Sub(h.embedFS, "build")
		if err != nil {
			return err
		}

		httpFS = http.FS(f)
	}
	router.NotFound = http.FileServer(httpFS)

	h.httpServer = &http.Server{
		Addr:    addr,
		Handler: router,
	}

	return nil
}

func (h *Homed) httpRender(w http.ResponseWriter, status string, data interface{}) {
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

func (h *Homed) httpRenderJSON(w http.ResponseWriter, data interface{}) {
	h.httpRender(w, "success", data)
}

func (h *Homed) httpError(w http.ResponseWriter, msg string) {
	h.httpRender(w, "error", msg)
}

func (h *Homed) httpComponentList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	h.httpRenderJSON(w, h.components)
}

func (h *Homed) updateComponent(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

func (h *Homed) websocketEvents(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
