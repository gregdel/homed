package apps

import (
	"context"
	"fmt"
	"sync"

	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/config"
	"go.uber.org/zap"
)

var applications map[string]App

func init() {
	applications = map[string]App{}
}

// App reprensents an app
type App interface {
	Name() string
	Init(*config.Config) error
	Run(*RunCtx) error
}

// RunCtx reprensents the run context of an app
type RunCtx struct {
	Ctx              context.Context
	Logger           *zap.Logger
	Components       *components.Components
	ComponentUpdated chan components.Component
}

// Register registers a new app
func Register(app App) {
	name := app.Name()
	_, ok := applications[name]
	if ok {
		panic(fmt.Errorf("apps: app %s already regitered", name))
	}

	applications[name] = app
}

// Apps reprensents the apps
type Apps struct {
	wg sync.WaitGroup
	m  map[string]App
}

// New returns a new app handler
func New() *Apps {
	return &Apps{
		m:  applications,
		wg: sync.WaitGroup{},
	}
}

// Init inits the apps and must be called before running the apps
func (a *Apps) Init(config *config.Config) error {
	for _, app := range a.m {
		if err := app.Init(config); err != nil {
			return err
		}
	}

	return nil
}

// Run starts all the apps
func (a *Apps) Run(ctx *RunCtx) {
	for _, app := range a.m {
		a.wg.Add(1)
		go func(app App) {
			defer a.wg.Done()

			z := zap.String("app_name", app.Name())

			ctx.Logger.Info("starting app", z)
			if err := app.Run(ctx); err != nil {
				ctx.Logger.Error("run failed", z, zap.Error(err))
				return
			}
			ctx.Logger.Info("app stopped", z)
		}(app)
	}

	a.wg.Wait()
}
