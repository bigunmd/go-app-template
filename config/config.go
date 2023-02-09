package config

import (
	"strings"

	"github.com/spf13/viper"
)

const ConfigFileName string = ".app.config"
const ConfigFileType string = "yaml"

func LoadConfig(configFilePath string) error {
	viper.AddConfigPath("~/")
	viper.AddConfigPath("./")
	viper.SetConfigName(ConfigFileName)
	viper.SetConfigType(ConfigFileType)
	if !strings.EqualFold(configFilePath, "") {
		viper.SetConfigFile(configFilePath)
	}
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
}

func SetDefaults() {
	HttpConfigSetDefault()
	FiberConfigSetDefault()
	PostgresConfigSetDefault()
	RedisConfigSetDefault()
}
