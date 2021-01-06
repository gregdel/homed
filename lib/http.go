package homed

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (h *Homed) initHTTP(addr string) error {
	router := httprouter.New()
	router.Handler("GET", "/metrics", promhttp.Handler())
	h.httpServer = &http.Server{
		Addr:    addr,
		Handler: router,
	}

	return nil
}
