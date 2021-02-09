package main

import (
	"github.com/evilsocket/shieldwall/api"
	"github.com/evilsocket/shieldwall/database"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Api      api.Config      `yaml:"api"`
	Email    api.EmailConfig `yaml:"mail"`
	Database database.Config `yaml:"database"`
}

func LoadConfig(path string) (*Config, error) {
	config := Config{}

	if err := cleanenv.ReadConfig(path, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
