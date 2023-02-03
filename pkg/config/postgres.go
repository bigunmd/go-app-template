package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

var Postgres postgres

type postgres struct {
	Host            string        `default:"127.0.0.1" json:"host,omitempty"`
	Port            int           `default:"5432" json:"port,omitempty"`
	SslMode         string        `default:"disable" split_words:"true" json:"ssl_mode,omitempty"`
	User            string        `default:"postgres" json:"user,omitempty"`
	Password        string        `default:"postgres" json:"password,omitempty"`
	Db              string        `default:"postgres" json:"db,omitempty"`
	ConnMaxIdleTime time.Duration `default:"10m" split_words:"true" json:"conn_max_idle_time,omitempty"`
	ConnMaxLifetime time.Duration `default:"1h" split_words:"true" json:"conn_max_lifetime,omitempty"`
	MaxIdleConns    int           `default:"2" split_words:"true" json:"max_idle_conns,omitempty"`
	MaxOpenConns    int           `default:"0" split_words:"true" json:"max_open_conns,omitempty"`
}

func LoadPostgresConfig() error {
	err := envconfig.Process("postgres", &Postgres)
	if err != nil {
		return err
	}
	return nil
}

func (p *postgres) GetConfigString() string {
	return fmt.Sprintf(
		`user=%s 
		password=%s 
		sslmode=%s 
		dbname=%s 
		host=%s 
		port=%v`,
		p.User,
		p.Password,
		p.SslMode,
		p.Db,
		p.Host,
		p.Port,
	)
}

func (p *postgres) GetUrl() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%v/%s?sslmode=%s",
		p.User,
		p.Password,
		p.Host,
		p.Port,
		p.Db,
		p.SslMode,
	)
}
