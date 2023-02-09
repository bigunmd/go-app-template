package db

import (
	"app/config"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

func NewRedisClient() *redis.Client {
	opt := new(redis.Options)
	opt.Addr = config.RedisAddr()
	opt.DB = viper.GetInt("redis.database")
	opt.PoolSize = viper.GetInt("redis.poolSize")
	if viper.GetBool("redis.auth") {
		opt.Username = viper.GetString("redis.username")
		opt.Password = viper.GetString("redis.password")
	}
	return redis.NewClient(opt)
}
