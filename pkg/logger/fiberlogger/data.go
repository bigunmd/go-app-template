package fiberlogger

import "time"

// data is a struct to define some variables to use in custom logger function
type data struct {
	pid   int
	start time.Time
	end   time.Time
}