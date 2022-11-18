package homed

import (
	"os"

	yaml "gopkg.in/yaml.v2"
)

func readFile(path string, data interface{}) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return yaml.NewDecoder(file).Decode(data)
}
