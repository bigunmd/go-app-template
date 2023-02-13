package main

import (
	"app/config"
	"app/infrastructure/db"
	"app/pkg/logger"
	"app/pkg/migrator"
	"app/service"
	"embed"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

//go:embed infrastructure/migrations
var migrations embed.FS

//	@title			Swagger Example API
//	@version		1.0
//	@description	This is a sample service template.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		127.0.0.1:8000
// @BasePath	/api
func main() {
	// Parse cli args to config
	pflag.String("config.filePath", ".app.config.yaml", "configuration file")
	pflag.String("http.host", "0.0.0.0", "serve http on specified host")
	pflag.Int("http.port", 8000, "serve http on specified port")
	pflag.String("logger.level", "info", "set logger level [panic, fatal, error, warn, info, debug, trace]")
	pflag.Bool("logger.writeToFile", false, "write logs to filesystem")
	pflag.String("logger.file.path", "~/", "where to store log files")
	pflag.String("logger.file.name", "app.log", "log file base name")
	pflag.String("logger.file.maxAge", "24h", "log file max age")
	pflag.String("logger.file.rotationTime", "168h", "log file rotation time")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)

	// Load configuration defaults
	logger.LoggerConfigSetDefault()
	config.SetDefaults()
	// Load config
	config.LoadConfig()
	// Create logger instance
	log := logger.NewLogger()
	// Load migration source from embeded filesystem
	source, err := iofs.New(migrations, "infrastructure/migrations")
	if err != nil {
		log.AddError(err).Fatal("cannot make source driver")
	}
	// Create migrator instance and apply migrations
	m, err := migrator.NewPgMigrator(config.PostgresUrl(), source, log)
	if err != nil {
		log.AddError(err).Fatal("cannot create postgres migrator")
	}
	err = m.MigrateUP()
	if err != nil {
		if err == migrate.ErrNoChange {
			log.AddError(err).Debug("cannot migrate up")
		} else {
			log.AddError(err).Fatal("cannot migrate up")
		}
	}
	// Create postgres database instance
	pg, err := db.NewPostgres()
	if err != nil {
		log.AddError(err).Fatal("cannot create postgres connection")
	}
	defer pg.Close()

	// Create redis client instance
	rc := db.NewRedisClient()
	defer rc.Close()

	// Create http service instance
	s := service.NewHTTPService(log, pg, rc)
	s.RegisterUtilityRoutes()
	s.RegisterUserRoutes()
	s.RegisterNotFoundRoutes()

	// Start http service on in a separate goroutine instance
	go func() {
		log.AddError(s.Serve()).Error()
	}()

	// Create channel to signify a signal being sent
	c := make(chan os.Signal, 1)
	// When an interrupt or termination signal is sent, notify the channel
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// This blocks the main thread until an interrupt is received
	<-c
	log.Info("gracefully shutting down...")
	_ = s.Shutdown()
	log.Info("running cleanup tasks...")
	log.Info("service shutdown successfully")
}
