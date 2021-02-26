package mqttd

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gregdel/homed/lib/apps"
	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/config"
	"go.uber.org/zap"
)

const name = "mqttd"

func init() {
	apps.Register(app())
}

type mqttd struct {
	logger     *zap.Logger
	components *components.Components
	config     *config.Config
	updateChan chan components.Component

	client      mqtt.Client
	stateTopics map[string]components.Component
	cmdTopics   map[string]components.Component
}

func app() *mqttd {
	return &mqttd{
		stateTopics: map[string]components.Component{},
		cmdTopics:   map[string]components.Component{},
	}
}

func (m *mqttd) Name() string {
	return name
}

func (m *mqttd) Init(config *config.Config) error {
	m.config = config
	return nil
}

func (m *mqttd) Run(ctx *apps.RunCtx) error {
	m.logger = ctx.Logger
	m.components = ctx.Components
	m.updateChan = ctx.ComponentUpdated

	for _, c := range ctx.Components.List() {
		if m.client == nil && c.MQTTClient() != nil {
			m.client = c.MQTTClient()
		}

		cfg := c.Config()
		if c.Internal() {
			m.cmdTopics[cfg.CommandTopic] = c
		}

		m.stateTopics[cfg.StateTopic] = c
	}

	m.logger.Info("connecting to MQTT")
	token := m.client.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	for topic := range m.stateTopics {
		m.logger.Info("subscribing to status topic", zap.String("topic", topic))
		token = m.client.Subscribe(topic, 0, m.handleMessage)
		if token.Wait() && token.Error() != nil {
			return token.Error()
		}
	}

	for topic := range m.cmdTopics {
		m.logger.Info("subscribing to command topic", zap.String("topic", topic))
		token = m.client.Subscribe(topic, 0, m.handleCommand)
		if token.Wait() && token.Error() != nil {
			return token.Error()
		}
	}

	<-ctx.Ctx.Done()

	m.logger.Info("disconnecting from the MQTT broker")
	m.client.Disconnect(250)

	return nil
}

func (m *mqttd) handleMessage(c mqtt.Client, msg mqtt.Message) {
	component, ok := m.stateTopics[msg.Topic()]
	if !ok {
		m.logger.Warn("Topic not found", zap.String("topic", msg.Topic()))
		return
	}

	if err := component.Update(msg.Payload()); err != nil {
		m.logger.Warn(
			"failed to update component",
			zap.Error(err),
			zap.String("friendly_name", string(component.FriendlyName())),
			zap.String("room", component.Room()),
			zap.String("device", component.Device()))
		return
	}

	if err := component.PostUpdate(); err != nil {
		m.logger.Warn(
			"failed to run the component post update",
			zap.Error(err),
			zap.String("friendly_name", string(component.FriendlyName())),
			zap.String("room", component.Room()),
			zap.String("device", component.Device()))
		return
	}

	m.updateChan <- component
}

func (m *mqttd) handleCommand(c mqtt.Client, msg mqtt.Message) {
	component, ok := m.cmdTopics[msg.Topic()]
	if !ok {
		m.logger.Warn("Topic not found", zap.String("topic", msg.Topic()))
		return
	}

	if err := component.ExecCommand(msg.Payload()); err != nil {
		m.logger.Warn(
			"failed to write component command",
			zap.String("error", err.Error()))
		return
	}

	m.logger.Debug(
		"Writing component command",
		zap.String("topic", msg.Topic()),
		zap.String("value", string(msg.Payload())),
	)
}
