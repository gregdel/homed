package homed

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (h *Homed) initHTTP(addr string) error {
	router := httprouter.New()
	router.Handler("GET", "/metrics", promhttp.Handler())
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
