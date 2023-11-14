package mqttd

import (
	"context"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gregdel/homed/lib/apps"
	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/config"
	"go.uber.org/zap"
)

func init() {
	apps.Register(app())
}

type mqttd struct {
	logger     *zap.Logger
	components *components.Components
	config     *config.Config

	mu          sync.Mutex
	client      mqtt.Client
	errChan     chan error
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
	return "mqttd"
}

func (m *mqttd) Init(config *config.Config) error {
	m.config = config

	opts := mqtt.
		NewClientOptions().
		AddBroker(config.MQTT.Broker).
		SetOnConnectHandler(m.onConnectHandler).
		SetReconnectingHandler(m.reconnectingHandler).
		SetConnectionLostHandler(m.connectionLostHandler)
	m.client = mqtt.NewClient(opts)

	return nil
}

func (m *mqttd) onConnectHandler(mqtt.Client) {
	m.logger.Info("connected to the broker")

	topics := map[string]byte{}
	for topic := range m.stateTopics {
		if topic != "" {
			topics[topic] = 0
		}
	}

	m.logger.Info("subscribing to state topics")
	token := m.client.SubscribeMultiple(topics, m.handleMessage)
	if token.Wait() && token.Error() != nil {
		m.logger.Error("failed to subscribe to the state topics",
			zap.Error(token.Error()),
		)
		m.errChan <- token.Error()
		return
	}

	topics = map[string]byte{}
	for topic := range m.cmdTopics {
		if topic != "" {
			topics[topic] = 0
		}
	}

	m.logger.Info("subscribing to command topics")
	token = m.client.SubscribeMultiple(topics, m.handleCommand)
	if token.Wait() && token.Error() != nil {
		m.logger.Error("failed to subscribe to the command topics",
			zap.Error(token.Error()),
		)
		m.errChan <- token.Error()
		return
	}
}

func (m *mqttd) reconnectingHandler(mqtt.Client, *mqtt.ClientOptions) {
	m.logger.Info("attempting to connect to the mqtt broker")
}

func (m *mqttd) connectionLostHandler(mqtt.Client, error) {
	m.logger.Info("connection to the mqtt broker is lost")
}

func (m *mqttd) Run(ctx context.Context, config *apps.Config) error {
	m.mu.Lock()
	m.logger = config.Logger.With(zap.String("app", m.Name()))
	m.components = config.Components
	m.errChan = make(chan error)

	for _, c := range m.components.List() {
		c.SetMQTTClient(m.client)

		cfg := c.Config()
		if c.Internal() {
			m.cmdTopics[cfg.CommandTopic] = c
		}

		m.stateTopics[cfg.StateTopic] = c
	}
	m.mu.Unlock()

	m.logger.Info("connecting to MQTT")
	token := m.client.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	var err error
	select {
	case <-ctx.Done():
		// Nothing to do
	case err = <-m.errChan:
		// Error
	}

	m.logger.Info("disconnecting from the MQTT broker")
	m.client.Disconnect(250)

	return err
}

func (m *mqttd) handleMessage(c mqtt.Client, msg mqtt.Message) {
	m.mu.Lock()
	component, ok := m.stateTopics[msg.Topic()]
	m.mu.Unlock()
	if !ok {
		m.logger.Warn("topic not found", zap.String("topic", msg.Topic()))
		return
	}

	logger := component.LoggerWithFields(m.logger)

	if err := component.Update(msg.Payload()); err != nil {
		logger.Warn("failed to update component", zap.Error(err))
		return
	}

	if err := component.PostUpdate(); err != nil {
		logger.Warn("failed to run the component post update", zap.Error(err))
		return
	}
}

func (m *mqttd) handleCommand(c mqtt.Client, msg mqtt.Message) {
	m.mu.Lock()
	component, ok := m.cmdTopics[msg.Topic()]
	m.mu.Unlock()
	if !ok {
		m.logger.Warn("topic not found", zap.String("topic", msg.Topic()))
		return
	}

	logger := component.LoggerWithFields(m.logger)

	if err := component.ExecCommand(msg.Payload()); err != nil {
		logger.Warn("failed to write component command", zap.Error(err))
		return
	}

	logger.Debug(
		"Writing component command",
		zap.String("topic", msg.Topic()),
		zap.String("value", string(msg.Payload())),
	)
}
