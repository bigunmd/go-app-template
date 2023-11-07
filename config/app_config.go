package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
)

type AppCfg struct {
	Logger   Logger   `json:"logger" yaml:"logger" env-prefix:"LOGGER_"`
	HTTP     HTTP     `json:"http" yaml:"http" env-prefix:"HTTP_"`
	TLS      TLS      `json:"tls" yaml:"tls" env-prefix:"TLS_"`
	Postgres Postgres `json:"postgres" yaml:"postgres" env-prefix:"POSTGRES_"`
	Redis    Redis    `json:"redis" yaml:"redis" env-prefix:"REDIS_"`
	Swagger  Swagger  `json:"swagger" yaml:"swagger" env-prefix:"SWAGGER_"`
}

func NewAppCfg(filepath string) (*AppCfg, error) {
	var err error
	var c AppCfg
	if filepath == "" {
		err = cleanenv.ReadEnv(&c)
		if err != nil {
			return nil, errors.Wrap(err, "cannot read env")
		}
	} else {
		err = cleanenv.ReadConfig(filepath, &c)
		if err != nil {
			return nil, errors.Wrap(err, "cannot read config")
		}
	}
	return &c, err
}
