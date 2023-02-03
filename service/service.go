package service

import (
	"app/pkg/config"
	"app/pkg/logger"
	"app/user"
	"app/utility"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"

	// "github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	fLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/helmet/v2"
	"github.com/gofiber/storage/redis"
	"github.com/jmoiron/sqlx"
)

type Service interface {
	Serve() error
	Shutdown() error
	RegisterUtilityRoutes()
	RegisterUserRoutes()
}

type service struct {
	log    logger.Logger
	router *fiber.App
	store  *redis.Storage
	db     *sqlx.DB
}

// RegisterUserRoutes implements Service
func (s *service) RegisterUserRoutes() {
	router := s.router.Group(config.Service.FullApiPath())
	c := user.NewUserController(router, s.db, s.log)
	c.RegisterRoutes()
}

// RegisterUtilityRoutes implements Service
func (s *service) RegisterUtilityRoutes() {
	router := s.router.Group("/")
	c := utility.NewUtilityController(router, s.log)
	c.RegisterRoutes()
}

// Shutdown implements Service
func (s *service) Shutdown() error {
	return s.router.Shutdown()
}

// Serve implements Service
func (s *service) Serve() error {
	return s.router.Listen(config.Service.GetUrl())
}

func NewService(logger logger.Logger, db *sqlx.DB) Service {
	s := redis.New(redis.Config{
		Host:     config.Redis.Host,
		Port:     config.Redis.Port,
		Database: config.Redis.Database,
		PoolSize: config.Redis.PoolSize,
	})
	f := fiber.New(config.Fiber)
	fLogConfig := fLogger.ConfigDefault
	fLogConfig.Format = "[${time}] ${pid} ${locals:requestid} ${ip}:${port} ${status} - ${method} ${url}\n"
	f.Use(
		recover.New(),
		fLogger.New(fLogConfig),
		compress.New(),
		cors.New(),
		// csrf.New(csrf.Config{
		// 	CookieSecure:   true,
		// 	CookieHTTPOnly: true,
		// 	Storage:        s,
		// 	ContextKey:     "csrf",
		// }),
		cache.New(cache.Config{
			Expiration:   config.FiberCache.Expiration,
			CacheControl: true,
			Storage:      s,
		}),
		etag.New(),
		encryptcookie.New(encryptcookie.Config{
			Key: config.FiberEncryptCookie.Key,
		}),
		favicon.New(favicon.Config{
			File: "./favicon.ico",
		}),
		limiter.New(limiter.Config{
			Next: func(c *fiber.Ctx) bool {
				return c.IP() == "127.0.0.1"
			},
			Max:               config.FiberLimiter.Max,
			Expiration:        config.FiberLimiter.Expiration,
			Storage:           s,
			LimiterMiddleware: limiter.SlidingWindow{},
		}),
		requestid.New(),
		helmet.New(),
	)
	return &service{
		log:    logger,
		router: f,
		store:  s,
		db:     db,
	}
}
