package router

import (
	"app/config"
	"app/pkg/logger"
	fiberlogger "app/pkg/middleware/fiberlogger"
	"embed"
	"net/http"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/helmet/v2"
	"github.com/gofiber/storage/redis"
	"github.com/spf13/viper"
)

//go:embed assets
var assets embed.FS

func NewFiberRouter(logger logger.Logger) *fiber.App {
	s := redis.New(config.NewFiberRedisStorageConfig())
	f := fiber.New(config.NewFiberConfig())
	f.Use(
		fiberlogger.New(logger),
		recover.New(),
		compress.New(),
		cors.New(config.NewFiberCorsConfig()),
		requestid.New(),
		helmet.New(),
		encryptcookie.New(
			encryptcookie.Config{Key: viper.GetString("fiber.encryptCookie.key")},
		),
		cache.New(config.NewFiberCacheConfig(s)),
		limiter.New(config.NewFiberLimiterConfig(s)),
		pprof.New(pprof.Config{
			Prefix: viper.GetString("fiber.pprof.prefix"),
		}),
	)
	if viper.GetBool("fiber.etag") {
		f.Use(etag.New())
	}
	if viper.GetBool("fiber.csrf.enable") {
		f.Use(csrf.New(csrf.Config{
			CookieSecure:   viper.GetBool("fiber.csrf.cookieSecure"),
			CookieHTTPOnly: viper.GetBool("fiber.csrf.cookieHttpOnly"),
			Storage:        s,
			ContextKey:     "csrf",
		}))
	}
	f.Use(
		favicon.New(favicon.Config{
			File:       "assets/favicon.ico",
			URL:        "/favicon.ico",
			FileSystem: http.FS(assets),
		}),
	)
	prometheus := fiberprometheus.New(viper.GetString("fiber.prometheus.serviceName"))
	prometheus.RegisterAt(f, viper.GetString("fiber.prometheus.path"))
	f.Use(prometheus.Middleware)
	return f
}
