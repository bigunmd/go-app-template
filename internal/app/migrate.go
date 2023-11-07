package app

import (
	"context"
	"goapptemplate"
	"goapptemplate/config"
	"goapptemplate/internal/domain"
	"goapptemplate/pkg/migrator"

	"github.com/pkg/errors"
)

func migrate(cfg *config.AppCfg) error {
	mu, err := migrator.NewPostgresMigrator(
		cfg.Postgres.ConfigURL(),
		domain.SchemaApp,
		goapptemplate.MigrationsApp,
		"migrations/app",
	)
	if err != nil {
		return errors.Wrap(err, "cannot create postgres migrator")
	}
	ctx, cancel := context.WithTimeout(context.Background(), cfg.HTTP.Timeout)
	defer cancel()
	err = mu.Up(ctx)
	if err != nil {
		return errors.Wrap(err, "cannot migrate up")
	}
	return nil
}
