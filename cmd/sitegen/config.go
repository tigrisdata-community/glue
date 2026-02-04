package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ContentDir string `yaml:"content_dir"`
	OutputDir  string `yaml:"output_dir"`
	Preamble   string `yaml:"preamble"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing YAML: %w", err)
	}

	if cfg.ContentDir == "" {
		return nil, fmt.Errorf("content_dir is required in config")
	}
	if cfg.OutputDir == "" {
		return nil, fmt.Errorf("output_dir is required in config")
	}

	cfg.ContentDir, _ = filepath.Abs(cfg.ContentDir)
	cfg.OutputDir, _ = filepath.Abs(cfg.OutputDir)

	return &cfg, nil
}
