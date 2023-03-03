package fakehome

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gregdel/homed/lib/apps"
	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/components/boiler"
	"github.com/gregdel/homed/lib/components/common"
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

	client    mqtt.Client
	cmdTopics map[string]components.Component

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
		SetConnectionLostHandler(fh.mqttOnConnectionLostHandler).
		SetDefaultPublishHandler(fh.commandHandler)
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
	fh.updateStates()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			fh.updateStates()
		}
	}
}

func (fh *FakeHome) mqttOnConnectHandler(c mqtt.Client) {
	fh.logger.Info("connected to mqtt")

	for topic := range fh.cmdTopics {
		fh.logger.Info("subscribing to command topic", zap.String("topic", topic))
		token := fh.client.Subscribe(topic, 0, nil)
		if token.Wait() && token.Error() != nil {
			fh.logger.Error("failed to subscribe to the command topic",
				zap.String("topic", topic),
				zap.Error(token.Error()),
			)
		}
	}
}

func (fh *FakeHome) mqttOnConnectionLostHandler(mqtt.Client, error) {
	fh.logger.Info("connection to the mqtt broker is lost")
}

func (fh *FakeHome) commandHandler(c mqtt.Client, msg mqtt.Message) {
	component, ok := fh.cmdTopics[msg.Topic()]
	if !ok {
		fh.logger.Warn("topic not found", zap.String("topic", msg.Topic()))
		return
	}

	payload := msg.Payload()

	var errUpdate error
	var errPublish error
	switch x := component.(type) {
	case *boiler.Boiler:
		errUpdate = x.Update(payload)
		errPublish = x.PublishToStateTopic(payload)
	case *common.RollerShutter:
		cancel, ok := fh.cancelFuncs[x.ID()]
		if ok {
			cancel()
		}

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
				break
			}

			ticker := time.NewTicker(500 * time.Millisecond)
			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					x.Value += 5 * factor
					if x.Value < 0 {
						x.Value = 0
					}
					if x.Value > 100 {
						x.Value = 100
					}

					valueStr := fmt.Sprintf("%.02f", x.Value)
					if err := x.PublishToStateTopic([]byte(valueStr)); err != nil {
						return
					}

					if x.Value == 0 || x.Value == 100 {
						return
					}
				}
			}
		}(ctx, payload)

		fh.cancelFuncs[x.ID()] = cancel
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
		fh.logger.Warn("failed to update component", zap.String("topic", msg.Topic()))
	}

	if errPublish != nil {
		fh.logger.Warn("failed to publish component state", zap.String("topic", msg.Topic()))
	}
}

func (fh *FakeHome) updateStates() {
	if !fh.client.IsConnectionOpen() {
		fh.logger.Info("mqtt broker not connected, not updating states")
		return
	}

	var err error
	for i, component := range fh.components.List() {
		switch c := component.(type) {
		case *zClimate.Sensor:
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

			c.Humidity = 60

			if c.Temp == 0 {
				c.Temp = 15
			} else if c.Temp < 10 {
				c.Temp = 10
			} else if c.Temp > 22 {
				c.Temp = 22
			} else {
				c.Temp = c.Temp + (factor * 0.1)
			}

			c.Pressure = 1000
			err = fh.publishStateJSON(c)
		case *trv.TRV:
			c.LocalTemperature = float64((i % 5) + 15)
			c.HeatingSetpoint = c.LocalTemperature + 1
			c.Mode = trv.SystemModeAuto
			if (i % 2) == 0 {
				c.LocalTemperatureCalibration = -1
				c.Position = 60
				c.Force = trv.ForceModeOpen
			}
			err = fh.publishStateJSON(c)
		case *common.PowerMeter:
			err = c.PublishToStateTopic([]byte(strconv.Itoa((i % 3 * 100))))
		case *common.DeviceStatus:
			err = c.PublishToStateTopic([]byte("online"))
		}

		if err != nil {
			fh.logger.Warn("failed to publish state", zap.Error(err))
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
