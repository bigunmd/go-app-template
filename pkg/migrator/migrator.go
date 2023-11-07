package migrator

import (
	"context"
	"embed"
	"fmt"
	"goapptemplate/pkg/postgres"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pkg/errors"
)

const createSchema = "create schema if not exists"
const dropSchema = "drop schema if exists"

func NewMigrationsSource(fs embed.FS, path string) (source.Driver, error) {
	// Load migration source from embeded filesystem
	src, err := iofs.New(fs, path)
	if err != nil {
		return nil, fmt.Errorf("cannot create source driver: %w", err)
	}
	return src, nil
}

type PostgresMigrator struct {
	m      *migrate.Migrate
	url    string
	schema string
}

// Down implements migrator.Migrator.
func (pm *PostgresMigrator) Down(ctx context.Context) error {
	err := pm.m.Down()
	if err != nil {
		return errors.Wrapf(err, "cannot migrate [%s] down", pm.schema)
	}
	return nil
}

// Up implements migrator.Migrator.
func (pm *PostgresMigrator) Up(ctx context.Context) error {
	err := pm.m.Up()
	if err != nil {
		return errors.Wrapf(err, "cannot migrate [%s] up", pm.schema)
	}
	return nil
}

func (pm *PostgresMigrator) CreateSchema(ctx context.Context, schema string) error {
	db, err := postgres.NewConn(ctx, pm.url, nil)
	if err != nil {
		return errors.Wrap(err, "cannot create connection")
	}
	defer db.Close(ctx)
	_, err = db.Exec(ctx, fmt.Sprintf("%s %s", createSchema, schema))
	if err != nil {
		return errors.Wrapf(err, "cannot execute create [%s] schema query", schema)
	}
	return nil
}

func (pm *PostgresMigrator) DropSchema(ctx context.Context, schema string) error {
	db, err := postgres.NewConn(ctx, pm.url, nil)
	if err != nil {
		return errors.Wrap(err, "cannot create connection")
	}
	defer db.Close(ctx)
	_, err = db.Exec(ctx, fmt.Sprintf("%s %s", dropSchema, schema))
	if err != nil {
		return errors.Wrapf(err, "cannot execute drop [%s] schema query", schema)
	}
	return nil
}

func NewPostgresMigrator(url string, schema string, migrations embed.FS, path string) (*PostgresMigrator, error) {
	pm := &PostgresMigrator{
		url:    url,
		schema: schema,
	}
	src, err := NewMigrationsSource(migrations, path)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create migrations source")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	err = pm.CreateSchema(ctx, pm.schema)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot create [%s] schema", pm.schema)
	}
	m, err := migrate.NewWithSourceInstance("iofs", src, fmt.Sprintf("%s&x-migrations-table=%s_migrations&search_path=%s", url, schema, schema))
	if err != nil {
		return nil, errors.Wrap(err, "cannot create migrate instance with source")
	}
	pm.m = m
	return pm, nil
}
