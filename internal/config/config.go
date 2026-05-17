package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	MediaDir        string
	DBPath          string
	Port            string
	MaxConcurrentDL int
	PollInterval    time.Duration
}

type yamlConfig struct {
	MediaDir        string `yaml:"media_dir"`
	DBPath          string `yaml:"db_path"`
	Port            string `yaml:"port"`
	MaxConcurrentDL int    `yaml:"max_concurrent_downloads"`
	PollInterval    string `yaml:"poll_interval"`
}

func defaults() Config {
	return Config{
		MediaDir:        "./data/media",
		DBPath:          "./data/ytdlp.db",
		Port:            "8080",
		MaxConcurrentDL: 2,
		PollInterval:    time.Hour,
	}
}

// Load reads config from configPath (YAML), then applies env var overrides.
// If configPath is empty, it defaults to ./config.yaml.
// A missing config file is not an error.
func Load(configPath string) (*Config, error) {
	if configPath == "" {
		configPath = "./config.yaml"
	}
	if v := os.Getenv("YTD_CONFIG_PATH"); v != "" {
		configPath = v
	}

	cfg := defaults()

	if err := loadYAML(configPath, &cfg); err != nil {
		return nil, err
	}

	if err := applyEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func loadYAML(path string, cfg *Config) error {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("reading config file: %w", err)
	}

	var yc yamlConfig
	if err := yaml.Unmarshal(data, &yc); err != nil {
		return fmt.Errorf("parsing config file: %w", err)
	}

	if yc.MediaDir != "" {
		cfg.MediaDir = yc.MediaDir
	}
	if yc.DBPath != "" {
		cfg.DBPath = yc.DBPath
	}
	if yc.Port != "" {
		cfg.Port = yc.Port
	}
	if yc.MaxConcurrentDL != 0 {
		cfg.MaxConcurrentDL = yc.MaxConcurrentDL
	}
	if yc.PollInterval != "" {
		d, err := time.ParseDuration(yc.PollInterval)
		if err != nil {
			return fmt.Errorf("config file: invalid poll_interval %q: %w", yc.PollInterval, err)
		}
		cfg.PollInterval = d
	}

	return nil
}

func applyEnv(cfg *Config) error {
	if v := os.Getenv("YTD_MEDIA_DIR"); v != "" {
		cfg.MediaDir = v
	}
	if v := os.Getenv("YTD_DB_PATH"); v != "" {
		cfg.DBPath = v
	}
	if v := os.Getenv("YTD_PORT"); v != "" {
		cfg.Port = v
	}
	if v := os.Getenv("YTD_MAX_CONCURRENT_DOWNLOADS"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 {
			return fmt.Errorf("invalid YTD_MAX_CONCURRENT_DOWNLOADS %q: must be a positive integer", v)
		}
		cfg.MaxConcurrentDL = n
	}
	if v := os.Getenv("YTD_POLL_INTERVAL"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return fmt.Errorf("invalid YTD_POLL_INTERVAL %q: %w", v, err)
		}
		cfg.PollInterval = d
	}

	return nil
}
