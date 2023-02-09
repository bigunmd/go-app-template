package db

import (
	"app/config"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

func NewPostgres() (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", config.PostgresConfigString())
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(viper.GetDuration("postgres.connMaxIdleTime"))
	db.SetConnMaxLifetime(viper.GetDuration("postgres.connMaxLifetime"))
	db.SetMaxIdleConns(viper.GetInt("postgres.maxIdleConn"))
	db.SetMaxOpenConns(viper.GetInt("postgres.maxOpenConn"))
	return db, nil
}
