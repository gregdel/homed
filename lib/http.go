package homed

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (h *Homed) initHTTP(addr string) error {
	router := httprouter.New()
	router.Handler("GET", "/metrics", promhttp.Handler())
	router.PUT("/components/:id", h.updateComponent)
	router.GET("/data", h.jsonData)
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

	err = component.WriteCommand(h.mqttClient, data)
	if err != nil {
		fmt.Fprintf(w, "failed to write mqtt command: %s", err.Error())
		return
	}
}
