package main

import "errors"

// Custom errors
var (
	ErrMissingFileName = errors.New("homectl: missing file name")
)
