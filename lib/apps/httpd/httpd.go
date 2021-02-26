package httpd

import (
	"context"
	"io/fs"
	"net/http"
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

const name = "httpd"

func init() {
	apps.Register(app())
}

type httpd struct {
	logger     *zap.Logger
	httpServer *http.Server
	components *components.Components
	render     *render.Render
	websockets map[string]*websocket.Conn
}

func app() *httpd {
	return &httpd{
		websockets: map[string]*websocket.Conn{},
		render:     render.New(),
	}
}

func (h *httpd) Init(config *config.Config) error {
	router := httprouter.New()
	router.Handler("GET", "/metrics", promhttp.Handler())
	router.GET("/events", h.websocketEvents)

	router.GET("/components", h.httpComponentList)
	router.PUT("/components/:id", h.updateComponent)
	router.GET("/components/:id/schedule", h.httpGetSchedule)
	router.POST("/components/:id/schedule/default", h.httpPostScheduleDefault)
	router.POST("/components/:id/schedule/overrides", h.httpPostScheduleOverrides)
	router.DELETE("/components/:id/schedule/overrides/:overrideID", h.httpDeleteScheduleOverride)
	router.POST("/components/:id/schedule/daily/:weekday", h.httpPostSchedule)
	router.DELETE("/components/:id/schedule/daily/:weekday/:uuid", h.httpDeleteSchedule)

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

	return nil
}

func (h *httpd) Name() string {
	return name
}

func (h *httpd) Run(ctx *apps.RunCtx) error {
	h.logger = ctx.Logger
	h.components = ctx.Components

	go func() {
		<-ctx.Ctx.Done()

		timeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		if err := h.httpServer.Shutdown(timeout); err != nil {
			h.logger.Error("failed to shutdown server", zap.Error(err))
		}
	}()

	go func() {
		for component := range ctx.ComponentUpdated {
			h.publishToWebsocket(component)
		}
	}()

	if err := h.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
}
