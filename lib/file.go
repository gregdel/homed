package homed

import (
	"path"
)

// FileType represents the configuration file types
type FileType string

// Custom file types
var (
	FileTypeRoom   FileType = "room"
	FileTypeSensor FileType = "sensor"
	FileTypeAction FileType = "action"
	FileTypeDevice FileType = "device"
)

// File represents a file
type File interface {
	FileName() string
	FileType() FileType

	configFormat() interface{}
	fromConfig(data interface{}, homed *Homed) error
	toConfig() (interface{}, error)
}

// NewFile returns a new file from a name and a filetype
func NewFile(name string, fileType FileType) (File, error) {
	var file File

	switch fileType {
	case FileTypeRoom:
		file = NewRoom(name)
	case FileTypeDevice:
		file = NewDevice(name)
	case FileTypeSensor:
		file = NewSensor(name)
	default:
		return nil, ErrInvalidFileType
	}

	return file, nil
}

// removeExt returns file path without the extension
func removeExt(filepath string) string {
	// Extension
	ext := path.Ext(filepath)
	// File length without the extension
	l := len(filepath) - len(ext)
	// Rebuild path
	return filepath[:l]
}
