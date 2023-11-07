package app

import (
	"context"
	"errors"
	"fmt"
	"goapptemplate/config"
	"goapptemplate/pkg/postgres"
	"os"
	"os/signal"
	"time"

	httpController "goapptemplate/internal/controller/http"
	"goapptemplate/internal/domain"
	"goapptemplate/internal/usecase"
	"goapptemplate/internal/usecase/repo"

	"github.com/gofiber/fiber/v2"
	"github.com/mikhail-bigun/fiberlogrus"

	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/helmet/v2"
	"github.com/gofiber/storage/redis"
	gomigrate "github.com/golang-migrate/migrate/v4"
	"github.com/sirupsen/logrus"
)

func Run(cfg *config.AppCfg) {
	// ________________________________________________________________________
	// Setup logger
	logger := newLogger(cfg)
	// ________________________________________________________________________
	// Migrate
	err := migrate(cfg)
	if err != nil {
		if errors.Is(err, gomigrate.ErrNoChange) {
			logger.WithError(err).Warning("cannot migrate")
		} else {
			logger.WithError(err).Fatal("cannot migrate")
		}
	} else {
		logger.Info("Successfully applied migrations")
	}
	// ________________________________________________________________________
	// Create Postgres database instance
	pgxTracer := postgres.NewLogrusQueryTracer(logger)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	db, err := postgres.NewPostgresDB(
		ctx,
		cfg.Postgres.ConfigString(
			fmt.Sprintf(
				"search_path=%s",
				domain.SchemaApp,
			),
		),
		pgxTracer,
	)
	if err != nil {
		logger.WithError(err).Fatal("cannot create postgres db")
	}
	// ________________________________________________________________________
	// Setup Fiber router
	f := fiber.New()
	// Add middleware
	f.Use(
		fiberlogrus.New(fiberlogrus.Config{
			Logger: logger,
			Tags: []string{
				fiberlogrus.TagLatency,
				fiberlogrus.TagMethod,
				fiberlogrus.TagURL,
				fiberlogrus.TagUA,
				fiberlogrus.TagBytesSent,
				fiberlogrus.TagPid,
				fiberlogrus.TagStatus,
			},
		}),
		recover.New(),
		compress.New(),
		cors.New(),
		helmet.New(),
		requestid.New(),
		etag.New(),
		pprof.New(),
		cache.New(cache.Config{
			CacheControl: true,
			Storage: redis.New(redis.Config{
				Host:     cfg.Redis.Host,
				Port:     int(cfg.Redis.Port),
				Username: cfg.Redis.Username,
				Password: cfg.Redis.Password,
				Database: cfg.Redis.DB,
			}),
			Next: func(c *fiber.Ctx) bool {
				return c.IP() == "127.0.0.1"
			},
		}),
	)
	// ________________________________________________________________________
	// Create Books repository
	br := repo.NewBooksPostgresRepo(db, logger)
	// Create Books usecase
	bu := usecase.NewBooks(br, logger)
	// Create App HTTP controller
	_ = httpController.NewAppHTTPController(
		f,
		bu,
		&httpController.AppHTTPControllerConfig{
			Timeout: cfg.HTTP.Timeout,
		},
		logger,
	)
	// ________________________________________________________________________
	// Not found handler last in stack
	f.Use(
		func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusNotFound).JSON(fiber.ErrNotFound)
		},
	)
	// Run Fiber router in a separate go routine
	go func() {
		err := runHTTP(f, cfg)
		if err != nil {
			logger.WithError(err).Fatal("cannot run HTTP")
		}
	}()
	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	// This blocks the main thread until an interrupt is received
	<-quit
	logger.Info("Gracefully shutting down...")
	err = f.Shutdown()
	if err != nil {
		logger.WithError(err).Fatal("cannot gracefully shutdown Fiber server")
	}
	logger.Info("Running cleanup tasks...")
	logger.Info("Service shutdown successfully")
}

func newLogger(cfg *config.AppCfg) *logrus.Logger {
	logger := logrus.New()
	lvl, err := logrus.ParseLevel(cfg.Logger.Level)
	if err != nil {
		logger.WithError(err).Fatalf("cannot parse logrus level [%s]", cfg.Logger.Level)
	}
	logger.SetLevel(lvl)
	logger.Formatter = &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "",
		// DisableSorting:         true,
		DisableLevelTruncation: true,
		PadLevelText:           true,
		QuoteEmptyFields:       true,
	}
	return logger
}

func runHTTP(f *fiber.App, cfg *config.AppCfg) error {
	if cfg.TLS.Cert.Filepath != "" &&
		cfg.TLS.Key.Filepath != "" {
		return f.ListenTLS(
			cfg.HTTP.Addr(),
			cfg.TLS.Cert.Filepath,
			cfg.TLS.Key.Filepath,
		)
	} else {
		return f.Listen(cfg.HTTP.Addr())
	}
}
