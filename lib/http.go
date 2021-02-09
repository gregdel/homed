package homed

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/schedule"
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
	router.POST("/components/:id/schedule/:weekday", h.httpPostSchedule)
	router.DELETE("/components/:id/schedule/:weekday/:uuid", h.httpDeleteSchedule)

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

func (h *Homed) httpError(w http.ResponseWriter, msg string, err error) {
	e := struct {
		Message string `json:"message"`
		Error   string `json:"error"`
	}{
		Message: msg,
		Error:   err.Error(),
	}

	h.render.JSON(w, http.StatusInternalServerError, e)
}

func (h *Homed) httpComponentList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := json.NewEncoder(w).Encode(h.components)
	if err != nil {
		fmt.Fprintf(w, "failed to encode data: %s", err.Error())
	}
}

func (h *Homed) updateComponent(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	component, err := h.components.Get(id)
	if err != nil {
		h.httpError(w, "failed to get the component", err)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.httpError(w, "failed to read body", err)
		return
	}

	err = component.WriteCommand(data)
	if err != nil {
		h.httpError(w, "failed to write mqtt command", err)
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
		h.httpError(w, "got error upgrading request", err)
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
		h.httpError(w, "failed to register websocket", err)
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

func (h *Homed) httpGetSchedule(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	componentID := ps.ByName("id")
	c, err := h.components.Get(componentID)
	if err != nil {
		fmt.Fprintf(w, "failed to get the component: %s", err.Error())
		return
	}

	sc, ok := c.(components.Scheduled)
	if !ok {
		fmt.Fprintf(w, "this component can not be scheduled")
		return
	}

	err = json.NewEncoder(w).Encode(sc.Schedule())
	if err != nil {
		fmt.Fprintf(w, "failed to encode data: %s", err.Error())
		return
	}
}

func (h *Homed) httpGetCompoment(ps httprouter.Params) (components.Component, error) {
	componentID := ps.ByName("id")
	return h.components.Get(componentID)
}

func (h *Homed) httpPostSchedule(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	componentID := ps.ByName("id")
	c, err := h.components.Get(componentID)
	if err != nil {
		fmt.Fprintf(w, "failed to get the component: %s", err.Error())
		return
	}

	sc, ok := c.(components.Scheduled)
	if !ok {
		fmt.Fprintf(w, "this component can not be scheduled")
		return
	}

	ts := schedule.TimeSlot{}
	err = json.NewDecoder(r.Body).Decode(&ts)
	if err != nil {
		fmt.Fprintf(w, "failed to decode data: %s", err.Error())
		return
	}

	weekdayStr := ps.ByName("weekday")
	weekday, err := strconv.Atoi(weekdayStr)
	if err != nil {
		fmt.Fprintf(w, "failed to parse weekday: %s", err.Error())
		return
	}

	if weekday < 0 || weekday > 6 {
		fmt.Fprintf(w, "invalid weekday")
		return
	}

	err = sc.ScheduleAdd(time.Weekday(weekday), &ts)
	if err != nil {
		fmt.Fprintf(w, "failed to add to the schedule: %s", err.Error())
		return
	}

	err = sc.SaveSchedule(h.components.SchedulePath(c))
	if err != nil {
		fmt.Fprintf(w, "failed to save schedule: %s", err.Error())
		return
	}
}

func (h *Homed) httpDeleteSchedule(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	componentID := ps.ByName("id")
	c, err := h.components.Get(componentID)
	if err != nil {
		fmt.Fprintf(w, "failed to get the component: %s", err.Error())
		return
	}

	sc, ok := c.(components.Scheduled)
	if !ok {
		fmt.Fprintf(w, "this component can not be scheduled")
		return
	}

	weekdayStr := ps.ByName("weekday")
	weekday, err := strconv.Atoi(weekdayStr)
	if err != nil {
		fmt.Fprintf(w, "failed to parse weekday: %s", err.Error())
		return
	}

	if weekday < 0 || weekday > 6 {
		fmt.Fprintf(w, "invalid weekday")
		return
	}

	uuid := ps.ByName("uuid")
	err = sc.ScheduleDelete(time.Weekday(weekday), uuid)
	if err != nil {
		fmt.Fprintf(w, "failed to add to the schedule: %s", err.Error())
		return
	}

	err = sc.SaveSchedule(h.components.SchedulePath(c))
	if err != nil {
		fmt.Fprintf(w, "failed to save schedule: %s", err.Error())
		return
	}
}
