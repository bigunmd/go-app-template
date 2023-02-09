package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func HttpConfigSetDefault() {
	_ = viper.BindEnv("http.host", "HTTP_HOST")
	viper.SetDefault("http.host", "0.0.0.0")
	_ = viper.BindEnv("http.port", "HTTP_PORT")
	viper.SetDefault("http.port", 8000)
	_ = viper.BindEnv("http.prefix", "HTTP_PREFIX")
	viper.SetDefault("http.prefix", "")
	_ = viper.BindEnv("http.apiPath", "HTTP_API_PATH")
	viper.SetDefault("http.apiPath", "/api")
}

func HttpListenAddr() string {
	return fmt.Sprintf("%s:%s", viper.GetString("http.host"), viper.GetString("http.port"))
}

func HttpFullApiPath() string {
	return fmt.Sprintf("%s%s", viper.GetString("http.prefix"), viper.GetString("http.apiPath"))
}
