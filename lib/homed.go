package homed

import (
	"context"
	"embed"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gregdel/homed/lib/apps"
	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/config"
)

// Homed needs to be used to load data efficiently
type Homed struct {
	config     *config.Config
	logger     *slog.Logger
	components *components.Components
}

// New returns a new Homed
func New(configPath string, embedFS *embed.FS) (*Homed, error) {
	homed := &Homed{}

	config := &config.Config{
		EmbedFS: embedFS,
	}
	if err := readFile(configPath, config); err != nil {
		return nil, err
	}
	homed.config = config

	homed.components = components.New(config.DataPath)

	level := slog.LevelInfo
	if config.Debug {
		level = slog.LevelDebug
	}
	homed.logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))

	for _, d := range config.Devices {
		for _, cfg := range d.Components {
			_, err := homed.components.Add(cfg, d.Room, d.Name)
			if err != nil {
				return nil, err
			}
		}
	}

	return homed, nil
}

// Run runs the app
func (h *Homed) Run() error {
	a := apps.New()
	if err := a.Init(h.config); err != nil {
		return err
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-sigs
		cancel()
	}()

	var wg sync.WaitGroup
	errCh := make(chan error, len(h.components.List()))
	for _, component := range h.components.List() {
		wg.Add(1)
		go func(c components.Component) {
			defer wg.Done()
			if err := c.Run(ctx, h.logger, h.components); err != nil {
				wrapped := fmt.Errorf("component %s (%s) failed: %w", c.ID(), c.Type(), err)
				h.logger.Error("component run failed",
					slog.String("component_id", c.ID()),
					slog.String("component_type", string(c.Type())),
					slog.Any("error", err),
				)
				select {
				case errCh <- wrapped:
				default:
				}
				cancel()
			}
		}(component)
	}

	runConfig := &apps.Config{
		Logger:     h.logger,
		Components: h.components,
	}

	// Wait for the apps
	a.Run(ctx, runConfig)
	cancel()

	// Wait for the components
	wg.Wait()

	select {
	case err := <-errCh:
		return err
	default:
		return nil
	}
}
