package components

import "errors"

// Custom errors
var (
	ErrNotImplemented      = errors.New("components: not implemented")
	ErrComponentNotFound   = errors.New("components: not found")
	ErrMissingDevice       = errors.New("components: missing device")
	ErrDeviceOffline       = errors.New("components: device is offline")
	ErrMissingMQTTClient   = errors.New("components: missing mqtt client")
	ErrComponentReadOnly   = errors.New("components: component is read only")
	ErrOperatingInProgress = errors.New("components: operation in progress")
	ErrExecNotInternal     = errors.New("components: only internal components have exec commands")
	ErrMissingScheduleName = errors.New("components: missing schedule name")
)
