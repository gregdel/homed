package fakehome

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gregdel/homed/lib/apps"
	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/components/boiler"
	"github.com/gregdel/homed/lib/components/common"
	rollershutter "github.com/gregdel/homed/lib/components/roller_shutter"
	zClimate "github.com/gregdel/homed/lib/components/zigbee2mqtt/climate_sensor"
	"github.com/gregdel/homed/lib/components/zigbee2mqtt/trv"
	"github.com/gregdel/homed/lib/config"
	"go.uber.org/zap"
)

func init() {
	apps.Register(newApp())
}

// FakeHome emulates the components of the home for development and testing
// purposes
type FakeHome struct {
	enabled bool

	logger     *zap.Logger
	components *components.Components
	config     *config.Config

	client      mqtt.Client
	mu          sync.RWMutex
	cmdTopics   map[string]components.Component
	cancelFuncs map[string]context.CancelFunc
}

func newApp() *FakeHome {
	return &FakeHome{
		cmdTopics:   map[string]components.Component{},
		cancelFuncs: map[string]context.CancelFunc{},
	}
}

// Name implements the App interface
func (fh *FakeHome) Name() string {
	return "fakehome"
}

// Init implements the App interface
func (fh *FakeHome) Init(c *config.Config) error {
	fh.enabled = c.FakeHome
	fh.config = c
	return nil
}

// Run implements the App interface
func (fh *FakeHome) Run(ctx context.Context, config *apps.Config) error {
	logger := config.Logger.With(zap.String("app", fh.Name()))

	if !fh.enabled {
		logger.Info("app is disabled")
		return nil
	}

	fh.logger = logger
	fh.components = config.Components

	opts := mqtt.NewClientOptions().
		AddBroker(fh.config.MQTT.Broker).
		SetOnConnectHandler(fh.mqttOnConnectHandler).
		SetConnectionLostHandler(fh.mqttOnConnectionLostHandler)
	fh.client = mqtt.NewClient(opts)

	for _, c := range config.Components.List() {
		cfg := c.Config()
		if !c.Internal() && cfg.CommandTopic != "" {
			fh.cmdTopics[cfg.CommandTopic] = c
		}
	}

	logger.Info("connecting to MQTT")
	token := fh.client.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	ticker := time.NewTicker(30 * time.Second)
	fh.updateLastSeen()
	fh.updateStates()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			fh.updateLastSeen()
			fh.updateStates()
		}
	}
}

func (fh *FakeHome) mqttOnConnectHandler(c mqtt.Client) {
	fh.logger.Info("connected to mqtt")

	fh.mu.RLock()
	topics := map[string]byte{}
	for topic := range fh.cmdTopics {
		topics[topic] = 0
	}
	fh.mu.RUnlock()

	token := fh.client.SubscribeMultiple(topics, fh.commandHandler)
	if token.Wait() && token.Error() != nil {
		fh.logger.Error("failed to subscribe to all the command topic",
			zap.Error(token.Error()),
		)
	}
}

func (fh *FakeHome) mqttOnConnectionLostHandler(mqtt.Client, error) {
	fh.logger.Info("connection to the mqtt broker is lost")
}

func (fh *FakeHome) commandHandler(c mqtt.Client, msg mqtt.Message) {
	fh.mu.RLock()
	component, ok := fh.cmdTopics[msg.Topic()]
	fh.mu.RUnlock()
	if !ok {
		fh.logger.Warn("topic not found", zap.String("topic", msg.Topic()))
		return
	}

	if !component.Device().IsOnline() {
		return
	}

	payload := msg.Payload()

	var errUpdate error
	var errPublish error
	switch x := component.(type) {
	case *boiler.Boiler:
		errUpdate = x.Update(payload)
		errPublish = x.PublishToStateTopic(payload)
	case *rollershutter.RollerShutter:
		fh.mu.RLock()
		cancel, ok := fh.cancelFuncs[x.ID()]
		if ok {
			cancel()
		}
		fh.mu.RUnlock()

		var ctx context.Context
		ctx, cancel = context.WithCancel(context.Background())
		go func(ctx context.Context, payload []byte) {
			factor := 1.0
			switch string(payload) {
			case "stop":
				return
			case "open":
				break
			case "close":
				factor = -1.0
			}

			ticker := time.NewTicker(500 * time.Millisecond)
			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					value := x.SensorValue()
					value += 5 * factor
					if value < 0 {
						value = 0
					}
					if value > 100 {
						value = 100
					}
					x.Value.Store(value)

					valueStr := fmt.Sprintf("%.02f", value)
					if err := x.PublishToStateTopic([]byte(valueStr)); err != nil {
						return
					}

					if value == 0 || value == 100 {
						return
					}
				}
			}
		}(ctx, payload)

		fh.mu.Lock()
		fh.cancelFuncs[x.ID()] = cancel
		fh.mu.Unlock()
	case *trv.TRV:
		errUpdate = x.Update(payload)
		errPublish = fh.publishStateJSON(x)
	case *common.BinaryLight:
		errUpdate = x.Update(payload)
		errPublish = x.PublishToStateTopic(payload)
	case *common.Switch:
		errUpdate = x.Update(payload)
		errPublish = x.PublishToStateTopic(payload)
	}

	if errUpdate != nil {
		fh.logger.Warn(
			"failed to update component",
			zap.String("topic", msg.Topic()),
			zap.Error(errUpdate))
	}

	if errPublish != nil {
		fh.logger.Warn(
			"failed to publish component state",
			zap.String("topic", msg.Topic()),
			zap.Error(errPublish),
		)
	}
}

func (fh *FakeHome) fakeSensor(c *zClimate.Sensor) {
	var isHeating = false
	if c.Dev != nil {
		tempController := fh.components.TemperatureController(c.Dev.Room)
		if tempController != nil {
			isHeating = tempController.IsHeating()
		}
	}

	var factor = -1.0
	if isHeating {
		factor = 1
	}

	humidity := c.HumidityV.Load()
	if humidity < 60 || humidity > 70 {
		humidity = 60
	} else {
		humidity += 2
	}
	c.HumidityV.Store(humidity)

	temp := c.TemperatureV.Load()
	if temp == 0 {
		temp = 15
	} else if temp < 10 {
		temp = 10
	} else if temp > 22 {
		temp = 22
	} else {
		temp = temp + (factor * 0.1)
	}
	c.TemperatureV.Store(temp)

	c.Pressure.Store(1000)
}

func (fh *FakeHome) updateLastSeen() {
	if !fh.client.IsConnectionOpen() {
		fh.logger.Info("mqtt broker not connected, not updating last seen")
		return
	}

	for _, component := range fh.components.List() {
		c, ok := component.(*common.DeviceStatus)
		if !ok {
			continue
		}

		err := c.PublishToStateTopic([]byte("online"))
		if err != nil {
			fh.logger.Warn("failed to set device online", zap.Error(err))
		}
	}
}

func (fh *FakeHome) updateStates() {
	if !fh.client.IsConnectionOpen() {
		fh.logger.Info("mqtt broker not connected, not updating states")
		return
	}

	var err error
	for i, component := range fh.components.List() {
		client := component.MQTTClient()
		if client == nil || !client.IsConnectionOpen() {
			fh.logger.Info("mqtt client is not connected, not updating states")
			return
		}

		switch c := component.(type) {
		case *zClimate.Sensor:
			fh.fakeSensor(c)
			err = fh.publishStateJSON(c)
		case *trv.TRV:
			c.LocalTemperature.Store(float64((i % 5) + 15))
			c.HeatingSetpoint.Store(c.LocalTemperature.Load() + 1)
			if i%2 == 0 {
				c.Force.Store(trv.ForceModeUnavailable)
			} else {
				c.Mode.Store(trv.SystemModeHeat)
			}
			if (i % 2) == 0 {
				c.LocalTemperatureCalibration.Store(-1)
				c.Position.Store(60)
				c.Force.Store(trv.ForceModeOpen)
			}
			err = fh.publishStateJSON(c)
		case *common.PowerMeter:
			err = c.PublishToStateTopic([]byte(strconv.Itoa((i % 3 * 100))))
		}

		if err != nil {
			fh.logger.Warn("failed to publish state",
				zap.String("function", "updateStates"),
				zap.Error(err))
		}

		err = nil
	}
}

func (fh *FakeHome) publishStateJSON(c components.Component) error {
	buf := bytes.Buffer{}
	if err := json.NewEncoder(&buf).Encode(c); err != nil {
		return err
	}
	return c.PublishToStateTopic(buf.Bytes())
}
