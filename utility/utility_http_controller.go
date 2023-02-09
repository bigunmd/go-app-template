package utility

import (
	_ "app/docs"
	"app/pkg/logger"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

type UtilityHttpController interface {
	RegisterRoutes()
}

type controller struct {
	log    logger.Logger
	router fiber.Router
}

// RegisterRoutes implements UtilityHttpController
func (c *controller) RegisterRoutes() {
	c.router.Get("/monitor", monitor.New())
	c.router.Get("/api/monitor", monitor.New(monitor.Config{APIOnly: true}))
	sr := c.router.Group("/swagger")
	sr.Get("*", swagger.HandlerDefault)
}

func NewUtilityController(router fiber.Router, logger logger.Logger) UtilityHttpController {
	return &controller{
		log:    logger,
		router: router,
	}
}
