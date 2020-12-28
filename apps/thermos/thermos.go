package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	homed "github.com/gregdel/homed/lib"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

// Thermos controls the temperature
type Thermos struct {
	Homed  *homed.Homed
	Boiler Boiler
	Mode   Mode
	Rooms  []*homed.Room

	httpServer *http.Server

	logger *zap.Logger

	sensorMap map[string]*homed.Sensor

	mqttClient mqtt.Client
}

// New returns a new thermos
func New(configPath, mqttBroker string, debug bool) *Thermos {
	opts := mqtt.NewClientOptions().AddBroker(mqttBroker)
	client := mqtt.NewClient(opts)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	httpServer := &http.Server{
		Addr:    ":9091",
		Handler: mux,
	}

	thermos := &Thermos{
		Homed:      homed.New(configPath),
		mqttClient: client,
		httpServer: httpServer,
	}

	var err error
	if debug {
		thermos.logger, err = zap.NewDevelopment()
	} else {
		thermos.logger, err = zap.NewProduction()
	}
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return thermos
}

func (t *Thermos) loadRooms() error {
	roomsNames, err := t.Homed.ListFiles(homed.FileTypeRoom)
	if err != nil {
		return err
	}

	t.Rooms = []*homed.Room{}
	for _, n := range roomsNames {
		r, err := t.Homed.Load(n, homed.FileTypeRoom)
		if err != nil {
			return err
		}

		t.Rooms = append(t.Rooms, r.(*homed.Room))
	}

	return nil
}

func (t *Thermos) handleMessage(c mqtt.Client, m mqtt.Message) {
	sensor, ok := t.sensorMap[m.Topic()]
	if !ok {
		t.logger.Warn("Topic not found", zap.String("topic", m.Topic()))
		return
	}

	if err := sensor.Update(string(m.Payload())); err != nil {
		t.logger.Warn(
			"failed to update sensor",
			zap.String("error", err.Error()))
		return

	}

	t.logger.Debug(
		"Updating sensor",
		zap.String("device", sensor.Device.Name),
		zap.String("sensor", sensor.Name),
		zap.Float64("value", sensor.Value),
	)
	t.logger.Sync()
}

func (t *Thermos) printState() {
	for _, room := range t.Rooms {
		for _, device := range room.Devices {
			for _, sensor := range device.Sensors {
				if sensor.Type != homed.SensorTypeTemperature {
					continue
				}

				fmt.Printf(
					"%s (%s): %f\n",
					room.Name,
					device.Name,
					sensor.Value,
				)
			}
		}
	}
}

// Run runs the app
func (t *Thermos) Run() error {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Map topics to sensors
	t.sensorMap = map[string]*homed.Sensor{}
	for _, room := range t.Rooms {
		for _, device := range room.Devices {
			for _, sensor := range device.Sensors {
				topic, err := sensor.MQTTTopic()
				if err != nil {
					return err
				}

				t.sensorMap[topic] = sensor
			}
		}
	}

	t.logger.Info("Connecting to MQTT")
	token := t.mqttClient.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	for topic := range t.sensorMap {
		t.logger.Info("Subscribing to topic", zap.String("topic", topic))
		token = t.mqttClient.Subscribe(topic, 0, t.handleMessage)
		if token.Wait() && token.Error() != nil {
			return token.Error()
		}
	}

	go func() {
		<-sigs
		t.httpServer.Shutdown(context.Background())
	}()

	t.logger.Info("Starting HTTP server")
	t.httpServer.ListenAndServe()

	t.logger.Info("Disconnecting from the MQTT broker")
	t.mqttClient.Disconnect(250)

	return nil
}
