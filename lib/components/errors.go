package components

import "errors"

// Custom errors
var (
	ErrNotImplemented    = errors.New("components: not implemented")
	ErrComponentNotFound = errors.New("components: not found")
)
