package homed

import (
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

// Config reprensents the configuration
type Config struct {
	Debug bool `yaml:"debug"`
	MQTT  struct {
		Broker string `yaml:"broker"`
	} `yaml:"mqtt"`
	HTTP struct {
		Addr string `yaml:"addr"`
	} `yaml:"http"`
	Devices []struct {
		Name    string `yaml:"name"`
		Room    string `yaml:"room"`
		Sensors []struct {
			Type  string `yaml:"type"`
			Topic string `yaml:"mqtt_topic"`
		} `yaml:"sensors"`
	} `yaml:"devices"`
}

func readFile(path string, data interface{}) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return yaml.NewDecoder(file).Decode(data)
}

func writeFile(path string, overwrite bool, data interface{}) error {
	if !overwrite {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return ErrFileExists
		}
	}

	// Create the missing directories
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	return yaml.NewEncoder(file).Encode(data)
}

func deleteFile(path string) error {
	return os.Remove(path)
}
