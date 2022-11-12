package fakehome

import (
	"bytes"
	"context"
	"encoding/json"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gregdel/homed/lib/apps"
	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/components/boiler"
	"github.com/gregdel/homed/lib/components/common"
	saswell "github.com/gregdel/homed/lib/components/saswell_trv"
	tuya "github.com/gregdel/homed/lib/components/tuya_trv"
	zClimate "github.com/gregdel/homed/lib/components/zigbee2mqtt/climate_sensor"
	"github.com/gregdel/homed/lib/config"
	"go.uber.org/zap"
)

const appName = "fakehome"

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
}

func newApp() *FakeHome {
	return &FakeHome{
		cmdTopics: map[string]components.Component{},
	}
}

// Name implements the App interface
func (fh *FakeHome) Name() string {
	return appName
}

// Init implements the App interface
func (fh *FakeHome) Init(c *config.Config) error {
	fh.enabled = c.FakeHome
	fh.config = c
	return nil
}

// Run implements the App interface
func (fh *FakeHome) Run(ctx context.Context, config *apps.Config) error {
	if !fh.enabled {
		config.Logger.Info("app is disabled", zap.String("app_name", appName))
		return nil
	}

	fh.logger = config.Logger
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

	fh.logger.Info("connecting to MQTT")
	token := fh.client.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	ticker := time.NewTicker(30 * time.Second)
	var exit bool
	fh.updateStates()
	for {
		select {
		case <-ctx.Done():
			exit = true
			return nil
		case <-ticker.C:
			fh.updateStates()
		}

		if exit {
			break
		}
	}

	return nil
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
	case *tuya.TuyaTRV:
		errUpdate = x.Update(payload)
		errPublish = fh.publishStateJSON(x)
	case *saswell.SaswellTRV:
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
			c.Humidity = 60
			c.Temp = float64((i % 3) + 15)
			if i%2 == 0 {
				c.Pressure = 1000
			}
			err = fh.publishStateJSON(c)
		case *tuya.TuyaTRV:
			c.LocalTemperature = float64((i % 5) + 15)
			err = fh.publishStateJSON(c)
		case *saswell.SaswellTRV:
			c.LocalTemperature = float64((i % 5) + 15)
			err = fh.publishStateJSON(c)
		case *common.PowerMeter:
			err = c.PublishToStateTopic([]byte(strconv.Itoa((i % 4 * 100))))
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
