package service

import (
	"app/config"
	"app/domain/service"
	"app/infrastructure/router"
	"app/pkg/logger"
	"app/user"
	"app/utility"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"

	"github.com/jmoiron/sqlx"
)

type svc struct {
	log    logger.Logger
	router *fiber.App
	redis  *redis.Client
	db     *sqlx.DB
}

// RegisterNotFoundRoutes implements service.HTTPService
func (s *svc) RegisterNotFoundRoutes() {
	s.router.Use(
		func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusNotFound).JSON(fiber.ErrNotFound)
		},
	)
}

// RegisterUserRoutes implements service.HTTPService
func (s *svc) RegisterUserRoutes() {
	router := s.router.Group(config.HttpFullApiPath())
	c := user.NewUserHTTPController(router, s.db, s.log)
	c.RegisterRoutes()
}

// RegisterUtilityRoutes implements service.HTTPService
func (s *svc) RegisterUtilityRoutes() {
	router := s.router.Group("/")
	c := utility.NewUtilityController(router, s.log)
	c.RegisterRoutes()
}

// Shutdown implements service.HTTPService
func (s *svc) Shutdown() error {
	return s.router.Shutdown()
}

// Serve implements service.HTTPService
func (s *svc) Serve() error {
	return s.router.Listen(config.HttpListenAddr())
}

func NewHTTPService(logger logger.Logger, db *sqlx.DB, rc *redis.Client) service.HTTPService {
	r := router.NewFiberRouter(logger)
	return &svc{
		log:    logger,
		router: r,
		redis:  rc,
		db:     db,
	}
}
