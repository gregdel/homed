package main

import "time"

type TimeRange struct {
	Start time.Time
	Stop  time.Time
}

// Schedule represents a schedule
type Schedule struct {
	Name     string
	HighTemp float64
	LowTemp  float64
}
