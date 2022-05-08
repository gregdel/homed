package homed

import (
	"context"
	"embed"
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gregdel/homed/lib/apps"
	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/config"
	"go.uber.org/zap"
)

// Homed needs to be used to load data efficiently
type Homed struct {
	config     *config.Config
	logger     *zap.Logger
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

	var err error
	if config.Debug {
		homed.logger, err = zap.NewDevelopment()
	} else {
		homed.logger, err = zap.NewProduction()
	}
	if err != nil {
		return nil, err
	}

	opts := mqtt.NewClientOptions().AddBroker(config.MQTT.Broker)
	mqttClient := mqtt.NewClient(opts)

	for _, d := range config.Devices {
		for _, cfg := range d.Components {
			_, err := homed.components.Add(cfg, mqttClient, d.Room, d.Name)
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

	componentChan := make(chan components.Component)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-sigs
		cancel()
		close(componentChan)
	}()

	runCtx := &apps.RunCtx{
		Ctx:              ctx,
		Logger:           h.logger,
		Components:       h.components,
		ComponentUpdated: componentChan,
	}

	a.Run(runCtx)

	return nil
}
