package config

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Config struct as defined
type Config struct {
	Clients   []ClientConfig `yaml:"clients"`
	SleepTime time.Duration  `yaml:"sleep_time"`
	LogLevel  string         `yaml:"log_level"`
}

type ClientConfig struct {
	Name       string     `yaml:"name"`
	APIKey     string     `yaml:"api_key"`
	Host       string     `yaml:"host"`
	Conditions Conditions `yaml:"conditions"`
	Options    Options    `yaml:"options"`
}

type Conditions struct {
	WaitingThreshold         time.Duration `yaml:"waiting_threshold"`
	DownloadTimeoutThreshold time.Duration `yaml:"download_timeout_threshold"`
	AverageSpeedThreshold    float64       `yaml:"average_speed_threshold"`
}

type Options struct {
	KeepInClient   bool `yaml:"keep_in_client"`
	BlockList      bool `yaml:"blocklist"`
	SkipRedownload bool `yaml:"skip_redownload"`
}

func LoadConfig(reader io.Reader) (*Config, error) {
	var config Config
	decoder := yaml.NewDecoder(reader)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("error decoding YAML: %v", err)
	}

	// Override APIKey from environment variables if available
	for i := range config.Clients {
		envVarName := "API_KEY_" + strings.ToUpper(config.Clients[i].Name)
		envAPIKey := os.Getenv(envVarName)
		if envAPIKey != "" {
			config.Clients[i].APIKey = envAPIKey
		}
	}

	return &config, nil
}
