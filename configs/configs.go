package configs

import (
	"gopkg.in/yaml.v2"
	"os"
)

var Version = "dev"

type Service struct {
	Name    string   `yaml:"name"`
	Command string   `yaml:"command"`
	Env     []string `yaml:"env"`
}

type Config struct {
	Services []Service `yaml:"services"`
}

func LoadConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
