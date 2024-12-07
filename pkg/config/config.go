package config

import (
	"fmt"
	"os"
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

func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("error decoding YAML: %v", err)
	}

	return &config, nil
}
