package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	homed "github.com/gregdel/homed/lib"
)

// Thermos controls the temperature
type Thermos struct {
	Homed  *homed.Homed
	Boiler Boiler
	Mode   Mode
	Rooms  []*homed.Room

	sensorMap map[string]*homed.Sensor

	mqttClient mqtt.Client
}

// New returns a new thermos
func New(configPath, mqttBroker string) *Thermos {
	opts := mqtt.NewClientOptions().AddBroker(mqttBroker)
	client := mqtt.NewClient(opts)

	return &Thermos{
		Homed:      homed.New(configPath),
		mqttClient: client,
	}
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
		fmt.Printf("Topic %s not found in sensorMap\n", m.Topic())
		return
	}

	sensor.Value = string(m.Payload())
	fmt.Printf(
		"Updating sensor %s:%s to %s\n",
		sensor.Device.Name,
		sensor.Name,
		sensor.Value,
	)
}

// Run runs the app
func (t *Thermos) Run() error {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

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

	fmt.Printf("Connecting to MQTT\n")
	token := t.mqttClient.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	for topic := range t.sensorMap {
		fmt.Printf("Subscibing to %s\n", topic)
		token = t.mqttClient.Subscribe(topic, 0, t.handleMessage)
		if token.Wait() && token.Error() != nil {
			return token.Error()
		}
	}

	<-sigs

	t.mqttClient.Disconnect(250)

	return nil
}
