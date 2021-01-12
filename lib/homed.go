package homed

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gregdel/homed/lib/components"
	"go.uber.org/zap"
)

// Homed needs to be used to load data efficiently
type Homed struct {
	basePath string

	logger *zap.Logger

	httpServer *http.Server

	mqttClient mqtt.Client

	rooms   map[string]*Room
	devices map[string]*Device

	topicComponents map[string]components.Component
}

// Rooms TODO delete
func (h *Homed) Rooms() map[string]*Room {
	return h.rooms
}

// New returns a new Homed
func New(configPath string) (*Homed, error) {
	homed := &Homed{
		rooms:   map[string]*Room{},
		devices: map[string]*Device{},

		topicComponents: map[string]components.Component{},
	}

	config := &Config{}
	if err := readFile(configPath, config); err != nil {
		return nil, err
	}

	var err error
	if config.Debug {
		homed.logger, err = zap.NewDevelopment()
	} else {
		homed.logger, err = zap.NewProduction()
	}
	if err != nil {
		return nil, err
	}

	for _, d := range config.Devices {
		_, ok := homed.devices[d.Name]
		if ok {
			return nil, ErrDuplicateDevice
		}

		room, ok := homed.rooms[d.Room]
		if !ok {
			room = NewRoom(d.Room)
			homed.rooms[d.Room] = room
		}

		device := NewDevice(d.Name)
		device.Room = room
		homed.devices[device.Name] = device
		room.AddDevice(device)

		for _, cfg := range d.Components {
			component, err := device.AddComponent(cfg)
			if err != nil {
				homed.logger.Warn(err.Error())
				continue
			}

			homed.topicComponents[cfg.StateTopic] = component
		}
	}

	opts := mqtt.NewClientOptions().AddBroker(config.MQTT.Broker)
	homed.mqttClient = mqtt.NewClient(opts)

	if err := homed.initHTTP(config.HTTP.Addr); err != nil {
		return nil, err
	}

	return homed, nil
}

func (h *Homed) handleMessage(c mqtt.Client, m mqtt.Message) {
	component, ok := h.topicComponents[m.Topic()]
	if !ok {
		h.logger.Warn("Topic not found", zap.String("topic", m.Topic()))
		return
	}

	if err := component.Update(m.Payload()); err != nil {
		h.logger.Warn(
			"failed to update component",
			zap.String("error", err.Error()))
		return
	}

	if err := component.PostUpdate(); err != nil {
		h.logger.Warn(
			"failed to run the component post update",
			zap.String("error", err.Error()))
		return
	}

	h.logger.Debug(
		"Updating component",
		zap.String("topic", m.Topic()),
		zap.String("value", string(m.Payload())),
	)
	h.logger.Sync()
}

// Run runs the app
func (h *Homed) Run() error {
	sigs := make(chan os.Signal, 1)
	done := make(chan struct{})
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	h.logger.Info("Connecting to MQTT")
	token := h.mqttClient.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	for topic := range h.topicComponents {
		h.logger.Info("Subscribing to topic", zap.String("topic", topic))
		token = h.mqttClient.Subscribe(topic, 0, h.handleMessage)
		if token.Wait() && token.Error() != nil {
			return token.Error()
		}
	}

	go func() {
		<-sigs
		close(done)
		h.httpServer.Shutdown(context.Background())
	}()

	// Start the temperature control function
	go h.temperatureControl(done)

	h.logger.Info("Starting HTTP server")
	h.httpServer.ListenAndServe()

	h.logger.Info("Disconnecting from the MQTT broker")
	h.mqttClient.Disconnect(250)

	return nil
}
