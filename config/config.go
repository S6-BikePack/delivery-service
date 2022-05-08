package config

import "github.com/spf13/viper"

type Config struct {
	Server      Server
	RabbitMQ    RabbitMQ
	Database    Database
	ServiceArea ServiceArea
}

type Server struct {
	Service string
	Port    string
}

type RabbitMQ struct {
	Host     string
	Port     int
	User     string
	Password string
	Exchange string
}

type Database struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	Debug    bool
}

type ServiceArea struct {
	Id         int
	Identifier string
}

func UseConfig(path string) (*Config, error) {
	v := viper.New()

	v.SetConfigName(path)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	v.Debug()

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config

	err := v.Unmarshal(&config)

	if err != nil {
		return nil, err
	}

	return &config, err
}
