package fiberlogger

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

// getLogrusFields calls FuncTag funcnions on matching keys
func getFields(ftm map[string]FuncTag, c *fiber.Ctx, d *data) map[string]interface{} {
	f := make(map[string]interface{})
	for k, ft := range ftm {
		f[k] = ft(c, d)
	}
	return f
}

// New creates a new middleware handler
func New(config Config) fiber.Handler {
	d := new(data)
	// Set PID once
	d.pid = os.Getpid()
	ftm := getFuncTagMap(config, d)

	return func(c *fiber.Ctx) error {
		d.start = time.Now()
		c.Next()
		d.end = time.Now()
		config.Logger.AddFields(getFields(ftm, c, d)).Info()
		return nil
	}
}
