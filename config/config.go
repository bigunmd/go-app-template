package config

import (
	"strings"

	"github.com/spf13/viper"
)

const ConfigFileName string = ".app.config"
const ConfigFileType string = "yaml"

func LoadConfig() {
	viper.AddConfigPath("~/")
	viper.AddConfigPath("./")
	viper.SetConfigName(ConfigFileName)
	viper.SetConfigType(ConfigFileType)
	if !strings.EqualFold(viper.GetString("config.filePath"), "") {
		viper.SetConfigFile(viper.GetString("config.filePath"))
	}
	viper.AutomaticEnv()
	_ = viper.ReadInConfig()
}

func SetDefaults() {
	HttpConfigSetDefault()
	FiberConfigSetDefault()
	PostgresConfigSetDefault()
	RedisConfigSetDefault()
}
