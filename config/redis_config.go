package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func RedisConfigSetDefault() {
	_ = viper.BindEnv("redis.auth", "REDIS_AUTH")
	viper.SetDefault("redis.auth", false)
	_ = viper.BindEnv("redis.host", "REDIS_HOST")
	viper.SetDefault("redis.host", "127.0.0.1")
	_ = viper.BindEnv("redis.port", "REDIS_PORT")
	viper.SetDefault("redis.port", 6379)
	_ = viper.BindEnv("redis.database", "REDIS_DATABASE")
	viper.SetDefault("redis.database", 0)
	_ = viper.BindEnv("redis.poolSize", "REDIS_POOL_SIZE")
	viper.SetDefault("redis.poolSize", 10)
	_ = viper.BindEnv("redis.username", "REDIS_USERNAME")
	viper.SetDefault("redis.username", "")
	_ = viper.BindEnv("redis.password", "REDIS_PASSWORD")
	viper.SetDefault("redis.password", "")
}

func RedisAddr() string {
	return fmt.Sprintf("%s:%s", viper.GetString("redis.host"), viper.GetString("redis.port"))
}
