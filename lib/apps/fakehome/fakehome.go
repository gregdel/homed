package fakehome

import (
	"bytes"
	"encoding/json"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gregdel/homed/lib/apps"
	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/components/boiler"
	"github.com/gregdel/homed/lib/components/common"
	status "github.com/gregdel/homed/lib/components/device_status"
	"github.com/gregdel/homed/lib/components/esphome"
	saswell "github.com/gregdel/homed/lib/components/saswell_trv"
	tuya "github.com/gregdel/homed/lib/components/tuya_trv"
	xiaomi "github.com/gregdel/homed/lib/components/xiaomi_aqara"
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
func (fh *FakeHome) Run(ctx *apps.RunCtx) error {
	if !fh.enabled {
		ctx.Logger.Info("app is disabled", zap.String("app_name", appName))
		return nil
	}

	fh.logger = ctx.Logger
	fh.components = ctx.Components

	opts := mqtt.NewClientOptions().AddBroker(fh.config.MQTT.Broker)
	fh.client = mqtt.NewClient(opts)

	for _, c := range ctx.Components.List() {
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

	for topic := range fh.cmdTopics {
		fh.logger.Info("subscribing to command topic", zap.String("topic", topic))
		token = fh.client.Subscribe(topic, 0, fh.handleCommand)
		if token.Wait() && token.Error() != nil {
			return token.Error()
		}
	}

	ticker := time.NewTicker(30 * time.Second)
	var exit bool
	fh.updateStates()
	for {
		select {
		case <-ctx.Ctx.Done():
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

func (fh *FakeHome) handleCommand(c mqtt.Client, msg mqtt.Message) {
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
		errPublish = fh.publishState(x, nil)
	case *tuya.TuyaTRV:
		errUpdate = x.Update(payload)
		errPublish = fh.publishState(x, nil)
	case *saswell.SaswellTRV:
		errUpdate = x.Update(payload)
		errPublish = fh.publishState(x, nil)
	case *esphome.Light:
		errUpdate = x.Update(payload)
		errPublish = fh.publishState(x, esphome.NewLightState(x.IsOn()))
	case *esphome.Switch:
		errUpdate = x.Update(payload)
		errPublish = fh.publishState(x, payload)
	}

	if errUpdate != nil {
		fh.logger.Warn("failed to update component", zap.String("topic", msg.Topic()))
	}

	if errPublish != nil {
		fh.logger.Warn("failed to publish component state", zap.String("topic", msg.Topic()))
	}
}

func (fh *FakeHome) updateStates() {
	var err error
	for _, component := range fh.components.List() {
		switch c := component.(type) {
		case *xiaomi.Climate:
			c.Humidity = 60
			c.Temp = 18
			c.Pressure = 1000
			err = fh.publishState(c, nil)
		case *tuya.TuyaTRV:
			c.LocalTemperature = 18
			err = fh.publishState(c, nil)
		case *saswell.SaswellTRV:
			c.LocalTemperature = 17
			err = fh.publishState(c, nil)
		case *common.PowerMeter:
			err = fh.publishState(c, []byte("150"))
		case *status.DeviceStatus:
			err = fh.publishState(c, []byte("online"))
		}

		if err != nil {
			fh.logger.Warn("failed to publish state", zap.Error(err))
		}

		err = nil
	}
}

func (fh *FakeHome) publishState(c components.Component, data interface{}) error {
	if data == nil {
		data = c
	} else {
		if b, ok := data.([]byte); ok {
			return c.PublishToStateTopic(b)
		}
	}

	buf := bytes.Buffer{}
	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		return err
	}

	return c.PublishToStateTopic(buf.Bytes())
}
