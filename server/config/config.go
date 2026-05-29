package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
	"zyrouge.me/umi/utils"
)

type Config struct {
	Server ServerConfig `toml:"server" validate:"required"`
}

type ServerConfig struct {
	Host        string   `toml:"host" validate:"required"`
	Port        int      `toml:"port" validate:"required"`
	CrossOrigin []string `toml:"cross_origin"`
	WebFiles    string   `toml:"web_files"`
}

var configCache *Config

func LoadConfig() error {
	file := os.Getenv("CONFIG_FILE")
	if file == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		file = filepath.Join(cwd, "config.toml")
	}
	bytes, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	config := Config{}
	err = toml.Unmarshal(bytes, &config)
	if err != nil {
		return err
	}
	err = config.Validate()
	if err != nil {
		return err
	}
	configCache = &config
	return nil
}

func (c Config) Validate() error {
	err := utils.GlobalValidator.Struct(c)
	if err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}
	return nil
}

func GetConfig() (*Config, error) {
	if configCache == nil {
		err := LoadConfig()
		if err != nil {
			return nil, err
		}
	}
	return configCache, nil
}
