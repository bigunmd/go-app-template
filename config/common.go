package config

import (
	"fmt"
	"time"
)

type Logger struct {
	Level string `json:"level" yaml:"level" env:"LEVEL" env-default:"info"`
}

type HTTP struct {
	Host    string        `json:"host" yaml:"host" env:"HOST" env-default:"0.0.0.0"`
	Port    int32         `json:"port" yaml:"port" env:"PORT" env-default:"8000"`
	Timeout time.Duration `json:"timeout" yaml:"timeout" env:"TIMEOUT" env-default:"4s"`
	Prefix  string        `json:"prefix" yaml:"prefix" env:"PREFIX" env-default:""`
	APIPath string        `json:"api_path" yaml:"apiPath" env:"API_PATH" env-default:"/api"`
}

func (http HTTP) Addr() string {
	return fmt.Sprintf("%s:%v", http.Host, http.Port)
}

func (http HTTP) FullAPIPath() string {
	return fmt.Sprintf("%s%s", http.Prefix, http.APIPath)
}

type TLS struct {
	Cert struct {
		Filepath string `json:"filepath" yaml:"filepath" env:"FILEPATH" env-default:""`
	} `json:"cert" yaml:"cert" env-prefix:"CERT_"`
	Key struct {
		Filepath string `json:"filepath" yaml:"filepath" env:"FILEPATH" env-default:""`
	} `json:"key" yaml:"key" env-prefix:"KEY_"`
}

type Postgres struct {
	Host              string        `json:"host" yaml:"host" env:"HOST" env-default:"127.0.0.1"`
	Port              int32         `json:"port" yaml:"port" env:"PORT" env-default:"5432"`
	SslMode           string        `json:"ssl_mode" yaml:"sslMode" env:"SSL_MODE" env-default:"disable"`
	Db                string        `json:"db" yaml:"db" env:"DB" env-default:"postgres"`
	User              string        `json:"user" yaml:"user" env:"USER" env-default:"postgres"`
	Password          string        `json:"password" yaml:"password" env:"PASSWORD" env-default:"postgres"`
	MaxConns          int           `json:"max_conns" yaml:"maxConns" env:"MAX_CONNS" env-default:"10"`
	MinConns          int           `json:"min_conns" yaml:"minConns" env:"MIN_CONNS" env-default:"2"`
	MaxConnLifetime   time.Duration `json:"max_conn_lifetime" yaml:"maxConnLifetime" env:"MAX_CONN_LIFETIME" env-default:"10m"`
	MaxConnIdleTime   time.Duration `json:"max_conn_idle_time" yaml:"maxConnIdleTime" env:"MAX_CONN_IDLE_TIME" env-default:"1m"`
	HealthCheckPeriod time.Duration `json:"health_check_period" yaml:"healthCheckPeriod" env:"HEALTH_CHECK_PERIOD" env-default:"10s"`
}

func (p Postgres) ConfigString(opts ...string) string {
	conf := fmt.Sprintf(
		"host=%s port=%v sslmode=%s user=%s password=%s dbname=%s pool_max_conns=%v pool_min_conns=%v pool_max_conn_lifetime=%s pool_max_conn_idle_time=%s pool_health_check_period=%s",
		p.Host,
		p.Port,
		p.SslMode,
		p.User,
		p.Password,
		p.Db,
		p.MaxConns,
		p.MinConns,
		p.MaxConnLifetime,
		p.MaxConnIdleTime,
		p.HealthCheckPeriod,
	)
	for _, v := range opts {
		conf += fmt.Sprintf(" %s", v)
	}
	return conf
}

func (p Postgres) ConfigURL(args ...string) string {
	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%v/%s?sslmode=%s",
		p.User,
		p.Password,
		p.Host,
		p.Port,
		p.Db,
		p.SslMode,
	)
	for _, v := range args {
		url += fmt.Sprintf("&%s", v)
	}
	return url
}

type Redis struct {
	Host     string `json:"host" yaml:"host" env:"HOST" env-default:"127.0.0.1"`
	Port     int32  `json:"port" yaml:"port" env:"PORT" env-default:"6379"`
	Username string `json:"username" yaml:"username" env:"USERNAME" env-default:""`
	Password string `json:"password" yaml:"password" env:"PASSWORD" env-default:""`
	DB       int    `json:"db" yaml:"db" env:"DB" env-default:"0"`
}

func (redis Redis) Addr() string {
	return fmt.Sprintf("%s:%v", redis.Host, redis.Port)
}

type Swagger struct {
	Host     string `json:"host" yaml:"host" env:"HOST" env-default:"127.0.0.1:8888"`
	BasePath string `json:"base_path" yaml:"basePath" env:"BASE_PATH" env-default:"/api"`
}
