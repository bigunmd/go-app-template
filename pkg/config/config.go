package config

import (
	"app/pkg/logger"
	"strings"

	"github.com/spf13/viper"
)

const ConfigFileName string = ".app.config"
const ConfigFileType string = "yaml"

func LoadConfig(path string) error {
	viper.AddConfigPath("~/")
	viper.AddConfigPath("./")
	viper.SetConfigName(ConfigFileName)
	viper.SetConfigType(ConfigFileType)
	if !strings.EqualFold(path, "") {
		viper.SetConfigFile(path)
	}
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	logger.LoggerConfigSetDefault()
	return nil
}
