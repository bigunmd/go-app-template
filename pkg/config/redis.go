package config

import "github.com/kelseyhightower/envconfig"

type redis struct {
	Host string `default:"127.0.0.1"`
	Port int    `default:"6379"`
	// Username string `default:"root"`
	// Password string `default:"root"`
	Database int `default:"0"`
	PoolSize int `default:"10" split_words:"true"`
}

var Redis redis

func LoadRedisConfig() error {
	err := envconfig.Process("redis", &Redis)
	if err != nil {
		return err
	}
	return nil
}
