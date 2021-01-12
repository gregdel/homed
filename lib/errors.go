package homed

import "errors"

// Custom errors
var (
	ErrFileExists          = errors.New("homed: file exists")
	ErrUnloadableFile      = errors.New("homed: unloadable file type")
	ErrInvalidFileType     = errors.New("homed: invalid file type")
	ErrInvalidConfigFormat = errors.New("homed: invalid config format")
	ErrMissingDevice       = errors.New("homed: missing device")
	ErrMissingRoom         = errors.New("homed: missing room")
	ErrDuplicateDevice     = errors.New("homed: duplicate device")
)
