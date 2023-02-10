package migrator

import (
	"app/pkg/logger"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type pgMigrator struct {
	m   *migrate.Migrate
	log logger.Logger
}

// MigrateDown implements Migrator
func (pgm *pgMigrator) MigrateDown() error {
	err := pgm.m.Down()
	if err != nil {
		return err
	}
	return nil
}

// MigrateUP implements Migrator
func (pgm *pgMigrator) MigrateUP() error {
	err := pgm.m.Up()
	if err != nil {
		return err
	}
	return nil
}

func NewPgMigrator(url string, source source.Driver, logger logger.Logger) (Migrator, error) {
	m, err := migrate.NewWithSourceInstance("iofs", source, url)
	if err != nil {
		return nil, err
	}
	return &pgMigrator{
		m:   m,
		log: logger,
	}, nil
}
