package config

import (
	"embed"

	"gopkg.in/yaml.v3"
)

// Location represents a location
type Location struct {
	Latitude  float64 `yaml:"latitude"`
	Longitude float64 `yaml:"longitude"`
	UTCOffset float64 `yaml:"utc_offset"`
}

// Config reprensents the configuration
type Config struct {
	Debug              bool               `yaml:"debug"`
	DataPath           string             `yaml:"data_path"`
	FakeHome           bool               `yaml:"fake_home"`
	Dev                bool               `yaml:"dev"`
	TemperatureControl TemperatureControl `yaml:"temperature_control"`
	RollerShutter      RollerShutter      `yaml:"roller_shutter"`
	Location           Location           `yaml:"location"`
	MQTT               struct {
		Broker string `yaml:"broker"`
	} `yaml:"mqtt"`
	HTTP struct {
		Addr string `yaml:"addr"`
	} `yaml:"http"`
	Devices []struct {
		Name       string      `yaml:"name"`
		Room       string      `yaml:"room"`
		Components []Component `yaml:"components"`
	} `yaml:"devices"`
	EmbedFS *embed.FS `yaml:"-"`
}

// Component represents the configuration of a component.
type Component struct {
	ID           string    `yaml:"id"`
	Type         string    `yaml:"type"`
	FriendlyName string    `yaml:"friendly_name"`
	StateTopic   string    `yaml:"state_topic"`
	CommandTopic string    `yaml:"command_topic"`
	Internal     bool      `yaml:"internal"`
	Hide         bool      `yaml:"hide"`
	ScheduleName string    `yaml:"schedule_name"`
	Params       yaml.Node `yaml:"params"`
	GraphURL     string    `yaml:"graph_url"`
}

// TemperatureControl represents the configuration of the temperature control
// daemon.
type TemperatureControl struct {
	Enabled              bool    `yaml:"enabled"`
	Hysteresis           float64 `yaml:"hysteresis"`
	CalibrateTRV         bool    `yaml:"calibrate_trv"`
	CalibrationThreshold float64 `yaml:"calibration_threshold"`
	CalibrationMaxOffset float64 `yaml:"calibration_max_offset"`
}

// RollerShutter represents the configuration for the automatic roller shutter
// daemon.
type RollerShutter struct {
	Enabled     bool `yaml:"enabled"`
	RandomDelay int  `yaml:"random_delay"`
}
