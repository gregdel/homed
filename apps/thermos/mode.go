package main

// Mode represents the possible thermos modes
type Mode string

var (
	// ModeCruiseControl keeps the temperatures as close as possible to the
	// given target
	ModeCruiseControl Mode = "cruse_control"
	// ModeScheduled follows a schedule to control the temperature
	ModeScheduled Mode = "scheduled"
	// ModeFreezePrevention keeps the house hot enough to prevent the pipes from freezing
	ModeFreezePrevention Mode = "freeze_prevention"
)
