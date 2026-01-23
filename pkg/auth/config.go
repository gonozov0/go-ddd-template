package auth

import (
	"errors"
	"fmt"

	"go-ddd-template/pkg/envutils"
)

const DefaultUserHeader = "X-User-Header"

type Config struct {
	UserHeader string
}

func (cfg *Config) validate() error {
	if cfg.UserHeader == "" {
		return errors.New("service ticket header must be set if service roles checking is enabled")
	}

	return nil
}

func LoadConfig() (Config, error) {
	var cfg Config

	if err := cfg.loadFromEnvs(); err != nil {
		return cfg, fmt.Errorf("could not load tvm config from envs: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return cfg, fmt.Errorf("tvm config validation failed: %w", err)
	}

	return cfg, nil
}

func (cfg *Config) loadFromEnvs() error {
	cfg.UserHeader = envutils.GetEnv("USER_HEADER", DefaultUserHeader)

	return nil
}
