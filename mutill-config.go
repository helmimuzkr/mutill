package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type MutillConfig struct {
	Command      string          `yaml:"command"`
	Services     []ServiceConfig `yaml:"services"`
	AutoShutdown bool            `yaml:"auto_shutdown"`
}

type ServiceConfig struct {
	Name string   `yaml:"name"`
	Path string   `yaml:"path"`
	Skip bool     `yaml:"skip"`
	Args []string `yaml:"args"`
}

func LoadConfig(path string) (*MutillConfig, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg MutillConfig
	if err := yaml.Unmarshal(f, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}
