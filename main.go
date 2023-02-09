package main

import (
	"app/config"
	"app/infrastructure/db"
	"app/pkg/logger"
	"app/service"
	"embed"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
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
	pflag.String("config.filePath", "", "configuration file")
	pflag.String("logger.level", "", "set logger level [panic, fatal, error, warn, info, debug, trace]")
	pflag.Bool("logger.writeToFile", false, "write logs to filesystem (default: false)")
	pflag.String("logger.file.path", "", "where to store log files (default: '~/')")
	pflag.String("logger.file.name", "", "log file base name (default: app.log)")
	pflag.String("logger.file.maxAge", "", "log file max age (default: 24h)")
	pflag.String("logger.file.rotationTime", "", "log file rotation time (default: 168h)")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)

	logger.LoggerConfigSetDefault()
	config.SetDefaults()
	config.LoadConfig()

	log := logger.NewLogger()
	source, err := iofs.New(migrations, "infrastructure/migrations")
	if err != nil {
		log.AddError(err).Fatal("cannot make source driver")
	}
	m, err := migrate.NewWithSourceInstance("iofs", source, config.PostgresUrl())
	if err != nil {
		log.AddError(err).Fatal("cannot source migrations")
	}
	err = m.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			log.AddError(err).Debug("cannot migrate up")
		} else {
			log.AddError(err).Fatal("cannot migrate up")
		}
	}
	pg, err := db.NewPostgres()
	if err != nil {
		log.AddError(err).Fatal("cannot create postgres connection")
	}
	defer pg.Close()

	rc := db.NewRedisClient()
	defer rc.Close()

	s := service.NewHTTPService(log, pg, rc)
	s.RegisterUtilityRoutes()
	s.RegisterUserRoutes()
	s.RegisterNotFoundRoutes()

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
