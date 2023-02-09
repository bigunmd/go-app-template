package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

func PostgresConfigSetDefault() {
	_ = viper.BindEnv("postgres.host", "POSTGRES_HOST")
	viper.SetDefault("postgres.host", "127.0.0.1")
	_ = viper.BindEnv("postgres.port", "POSTGRES_PORT")
	viper.SetDefault("postgres.port", 5432)
	_ = viper.BindEnv("postgres.sslMode", "POSTGRES_SSL_MODE")
	viper.SetDefault("postgres.sslMode", "disable")
	_ = viper.BindEnv("postgres.db", "POSTGRES_DB")
	viper.SetDefault("postgres.db", "postgres")
	_ = viper.BindEnv("postgres.user", "POSTGRES_USER")
	viper.SetDefault("postgres.user", "postgres")
	_ = viper.BindEnv("postgres.password", "POSTGRES_PASSWORD")
	viper.SetDefault("postgres.password", "postgres")
	_ = viper.BindEnv("postgres.maxOpenConn", "POSTGRES_MAX_OPEN_CONN")
	viper.SetDefault("postgres.maxOpenConn", 0)
	_ = viper.BindEnv("postgres.maxIdleConn", "POSTGRES_MAX_IDLE_CONN")
	viper.SetDefault("postgres.maxIdleConn", 2)
	_ = viper.BindEnv("postgres.connMaxIdleTime", "POSTGRES_CONN_MAX_IDLE_TIME")
	viper.SetDefault("postgres.connMaxIdleTime", 1*time.Minute)
	_ = viper.BindEnv("postgres.connMaxLifetime", "POSTGRES_CONN_MAX_LIFETIME")
	viper.SetDefault("postgres.connMaxLifetime", 1*time.Hour)
}

func PostgresConfigString() string {
	return fmt.Sprintf(
		`user=%s 
		password=%s 
		sslmode=%s 
		dbname=%s 
		host=%s 
		port=%s`,
		viper.GetString("postgres.user"),
		viper.GetString("postgres.password"),
		viper.GetString("postgres.sslMode"),
		viper.GetString("postgres.db"),
		viper.GetString("postgres.host"),
		viper.GetString("postgres.port"),
	)
}

func PostgresUrl() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		viper.GetString("postgres.user"),
		viper.GetString("postgres.password"),
		viper.GetString("postgres.host"),
		viper.GetString("postgres.port"),
		viper.GetString("postgres.db"),
		viper.GetString("postgres.sslMode"),
	)
}
