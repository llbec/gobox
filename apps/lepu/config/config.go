package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database DatabaseConfig `yaml:"database"`
	Server   ServerConfig   `yaml:"server"`
	Timezone string         `yaml:"timezone"`
}

type DatabaseConfig struct {
	URL      string `yaml:"url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

var TimezoneLocation *time.Location

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}

	// Set default timezone if not specified
	if cfg.Timezone == "" {
		cfg.Timezone = "Asia/Shanghai"
	}

	// Load timezone location
	TimezoneLocation, err = time.LoadLocation(cfg.Timezone)
	if err != nil {
		return nil, fmt.Errorf("failed to load timezone %s: %v", cfg.Timezone, err)
	}

	return &cfg, nil
}