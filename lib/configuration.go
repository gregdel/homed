package homed

import (
	"os"

	yaml "gopkg.in/yaml.v3"
)

func readFile(path string, data any) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	return yaml.NewDecoder(file).Decode(data)
}
