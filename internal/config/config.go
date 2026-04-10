package config

import (
	"encoding/json"
	"errors"
	"os"
)

var gConfig = Config{
	DetectType:   true,
	DefaultHosts: []string{},
}

type Config struct {
	DetectType   bool     `json:"detect_type"`
	DefaultHosts []string `json:"default_hosts"`
}

func Get() *Config {
	return &gConfig
}

func Load(path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if len(bytes) == 0 {
		return nil
	}

	err = json.Unmarshal(bytes, &gConfig)
	return err
}
