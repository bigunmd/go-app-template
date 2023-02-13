package fiberlogger

import (
	"app/pkg/logger"
)

// Config defines the config for middleware
type Config struct {
	Logger logger.Logger
	Tags   []string
}
