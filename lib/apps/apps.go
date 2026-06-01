package apps

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/config"
)

var applications map[string]App

func init() {
	applications = map[string]App{}
}

// App reprensents an app
type App interface {
	Name() string
	Init(*config.Config) error
	Run(context.Context, *Config) error
}

// Config reprensents the running config of an app
type Config struct {
	Logger     *slog.Logger
	Components *components.Components
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
func (a *Apps) Run(parentCtx context.Context, config *Config) {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	for _, app := range a.m {
		a.wg.Add(1)
		go func(app App, cancel context.CancelFunc) {
			defer a.wg.Done()

			appName := slog.String("app_name", app.Name())

			config.Logger.Info("starting app", appName)
			if err := app.Run(ctx, config); err != nil {
				config.Logger.Error("run failed", appName, slog.Any("error", err))
				cancel()
				return
			}
			config.Logger.Info("app stopped", appName)
		}(app, cancel)
	}

	a.wg.Wait()
}
