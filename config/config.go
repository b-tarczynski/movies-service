package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Api      Api
	Postgres Postgres
	Headers  Headers
	Tmdb     Tmdb
}

type Api struct {
	Port    string
	Release bool
}

type Postgres struct {
	Address  string
	User     string
	Password string
	Database string
}

type Tmdb struct {
	Url string
	Key string
}

type Headers struct {
	Account string
	Id      string
}

func NewConfig(fileName string) *Config {
	viper.SetConfigFile(fileName)
	viper.SetConfigType("yaml")

	var config Config

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	return &config
}
