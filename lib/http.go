package homed

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/gregdel/homed/lib/schedule"
	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (h *Homed) initHTTP(addr string) error {
	router := httprouter.New()
	router.Handler("GET", "/metrics", promhttp.Handler())
	router.PUT("/components/:id", h.updateComponent)
	router.GET("/data", h.jsonData)
	router.GET("/events", h.websocketEvents)

	router.GET("/schedules/:roomName", h.httpGetSchedule)
	router.POST("/schedules/:roomName/:weekday", h.httpPostSchedule)
	router.DELETE("/schedules/:roomName/:weekday/:uuid", h.httpDeleteSchedule)

	router.NotFound = http.FileServer(http.Dir("frontend/build"))
	h.httpServer = &http.Server{
		Addr:    addr,
		Handler: router,
	}

	return nil
}

func (h *Homed) jsonData(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	out := struct {
		Rooms []*Room `json:"rooms"`
	}{}

	for _, r := range h.rooms {
		out.Rooms = append(out.Rooms, r)
	}

	// TODO check stuff
	err := json.NewEncoder(w).Encode(out)
	if err != nil {
		fmt.Fprintf(w, "failed to encode data: %s", err.Error())
	}
}

func (h *Homed) updateComponent(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	component, ok := h.components[id]
	if !ok {
		fmt.Fprintf(w, "component %s not found", id)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "failed to read body: %s", err.Error())
		return
	}

	err = component.WriteCommand(data)
	if err != nil {
		fmt.Fprintf(w, "failed to write mqtt command: %s", err.Error())
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
		fmt.Fprintf(w, "got error upgrading request: %s", err.Error())
		if _, ok := err.(websocket.HandshakeError); !ok {
			fmt.Fprintf(w, "handshake error")
		}
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
		fmt.Fprintf(w, err.Error())
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
	// s, ok := h.temperatureController.schedules[ps.ByName("roomName")]
	// if !ok {
	// 	fmt.Fprintf(w, "failed to find schedule")
	// 	return
	// }

	// err := json.NewEncoder(w).Encode(s)
	// if err != nil {
	// 	fmt.Fprintf(w, "failed to encode data: %s", err.Error())
	// 	return
	// }
}

func (h *Homed) httpPostSchedule(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ts := schedule.TimeSlot{}
	err := json.NewDecoder(r.Body).Decode(&ts)
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

	// s, ok := h.temperatureController.schedules[ps.ByName("roomName")]
	// if !ok {
	// 	fmt.Fprintf(w, "failed to find schedule")
	// 	return
	// }

	// err = s.Add(time.Weekday(weekday), &ts)
	// if err != nil {
	// 	fmt.Fprintf(w, "failed to add to the schedule: %s", err.Error())
	// 	return
	// }

	// h.saveTemperatureSchedules()
}

func (h *Homed) httpDeleteSchedule(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	weekdayStr := ps.ByName("weekday")
	// uuid := ps.ByName("uuid")

	weekday, err := strconv.Atoi(weekdayStr)
	if err != nil {
		fmt.Fprintf(w, "failed to parse weekday: %s", err.Error())
		return
	}

	if weekday < 0 || weekday > 6 {
		fmt.Fprintf(w, "invalid weekday")
		return
	}

	// s, ok := h.temperatureController.schedules[ps.ByName("roomName")]
	// if !ok {
	// 	fmt.Fprintf(w, "failed to find schedule")
	// 	return
	// }

	// err = s.Delete(time.Weekday(weekday), uuid)
	// if err != nil {
	// 	fmt.Fprintf(w, "failed to add to the schedule: %s", err.Error())
	// 	return
	// }

	// h.saveTemperatureSchedules()
}
