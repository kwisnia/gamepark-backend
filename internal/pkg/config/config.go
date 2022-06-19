package config

import (
	"github.com/spf13/viper"
)

var Config *Configuration

type Configuration struct {
	Server   ServerConfiguration
	Database DatabaseConfiguration
}

type DatabaseConfiguration struct {
	Driver       string
	Dbname       string
	Username     string
	Password     string
	Host         string
	Port         string
	MaxLifetime  int
	MaxOpenConns int
	MaxIdleConns int
}

type ServerConfiguration struct {
	Port   string
	Secret string
	Mode   string
}

func LoadConfig() {
	var configuration *Configuration
	viper.SetConfigName("api-config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&configuration)
	if err != nil {
		panic(err)
	}

	Config = configuration
}