package config

import (
	"fmt"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type service struct {
	Host    string `default:"0.0.0.0" json:"host,omitempty"`
	Port    int    `default:"8000" json:"port,omitempty"`
	Prefix  string `default:"" json:"prefix,omitempty"`
	ApiPath string `default:"/api" split_words:"true" json:"api_path,omitempty"`
}

var Service service

func LoadServiceConfig(host string, port int) error {
	err := envconfig.Process("service", &Service)
	if err != nil {
		return err
	}
	if !strings.EqualFold(host, "") {
		Service.Host = host
	}
	if port != -1 {
		Service.Port = port
	}
	return nil
}

func (s *service) GetUrl() string {
	return fmt.Sprintf("%s:%v", s.Host, s.Port)
}

func (s *service) FullApiPath() string {
	return fmt.Sprintf("%s%s", s.Prefix, s.ApiPath)
}
