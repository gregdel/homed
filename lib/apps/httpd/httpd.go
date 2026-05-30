package httpd

import (
	"context"
	"io/fs"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/gregdel/homed/lib/apps"
	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/config"
	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/unrolled/render"
	"go.uber.org/zap"
)

func init() {
	apps.Register(app())
}

type httpd struct {
	logger        *zap.Logger
	httpServer    *http.Server
	httpClient    *http.Client
	components    *components.Components
	render        *render.Render
	prometheusURL string

	mu         sync.RWMutex
	websockets map[*websocket.Conn]string
	exiting    atomic.Bool
}

func app() *httpd {
	return &httpd{
		websockets: map[*websocket.Conn]string{},
		render:     render.New(),
	}
}

func (h *httpd) Init(config *config.Config) error {
	router := httprouter.New()
	router.Handler("GET", "/metrics", promhttp.Handler())
	router.GET("/events", h.websocketEvents)

	router.GET("/components", h.httpComponentList)
	router.PUT("/components/:id", h.updateComponent)
	router.GET("/components/:id/graph", h.getComponentGraph)
	router.GET("/components/:id/schedule", h.getSchedule)
	router.POST("/components/:id/schedule/default", h.updateScheduleDefault)
	router.POST("/components/:id/schedule/overrides", h.getScheduleOverrides)
	router.DELETE("/components/:id/schedule/overrides/:overrideID", h.deleteScheduleOverride)
	router.POST("/components/:id/schedule/daily/:weekday", h.addScheduleTimeSlot)
	router.DELETE("/components/:id/schedule/daily/:weekday/:tsID", h.deleteScheduleTimeSlot)
	router.PUT("/components/:id/schedule/daily/:weekday/:tsID", h.updateScheduleTimeSlot)

	var httpFS http.FileSystem
	if config.Dev {
		httpFS = http.Dir("frontend/build")
	} else {
		f, err := fs.Sub(config.EmbedFS, "build")
		if err != nil {
			return err
		}

		httpFS = http.FS(f)
	}
	router.NotFound = http.FileServer(httpFS)

	h.httpServer = &http.Server{
		Addr:    config.HTTP.Addr,
		Handler: router,
	}
	h.httpClient = &http.Client{Timeout: 10 * time.Second}
	h.prometheusURL = config.Prometheus.URL
	h.exiting.Store(false)

	return nil
}

func (h *httpd) Name() string {
	return "httpd"
}

func (h *httpd) Run(ctx context.Context, config *apps.Config) error {
	h.logger = config.Logger.With(zap.String("app", h.Name()))
	h.components = config.Components

	go func() {
		<-ctx.Done()
		h.exiting.Store(true)

		timeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		if err := h.httpServer.Shutdown(timeout); err != nil {
			h.logger.Error("failed to shutdown server", zap.Error(err))
		}
	}()

	go func() {
		events := make(chan components.Event)
		for _, c := range config.Components.List() {
			c.Subscribe(h.Name(), events)
		}

		for {
			select {
			case <-ctx.Done():
				return
			case event := <-events:
				h.publishToWebsocket(event.ID)
			}
		}
	}()

	if err := h.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
}
