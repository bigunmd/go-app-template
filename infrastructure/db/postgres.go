package db

import (
	"app/pkg/config"

	"github.com/jmoiron/sqlx"
)

func NewPostgres() (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", config.Postgres.GetConfigString())
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(config.Postgres.ConnMaxIdleTime)
	db.SetConnMaxLifetime(config.Postgres.ConnMaxLifetime)
	db.SetMaxIdleConns(config.Postgres.MaxIdleConns)
	db.SetMaxOpenConns(config.Postgres.MaxOpenConns)
	return db, nil
}
