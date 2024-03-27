package internal

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type Config struct {
	ServerPort        string
	InterruptTimeout  time.Duration
	SentryDSN         string
	SentryEnvironment string
	RedisClusterMode  bool
	RedisTLSEnabled   bool
	RedisAddress      string
	RedisUsername     string
	RedisPassword     string
	RedisExpiration   time.Duration
}

func LoadConfig() (Config, error) {
	var (
		config Config
		err    error
	)

	config.ServerPort = getEnv("SERVER_PORT", "8080")
	config.InterruptTimeout, err = time.ParseDuration(getEnv("KILL_TIMEOUT", "2s"))
	if err != nil {
		return config, fmt.Errorf("could not parse kill timeout: %w", err)
	}
	config.SentryEnvironment = getEnv("SENTRY_ENVIRONMENT", "development")
	config.SentryDSN = getEnv("SENTRY_DSN", "")
	config.RedisClusterMode = getEnv("REDIS_CLUSTER_MODE", "false") == "true"
	config.RedisTLSEnabled = getEnv("REDIS_TLS_ENABLED", "false") == "true"
	config.RedisAddress = getEnv("REDIS_ADDRESS", "localhost:6379")
	config.RedisUsername = getEnv("REDIS_USERNAME", "")
	config.RedisPassword = getEnv("REDIS_PASSWORD", "")
	redisExpiration := getEnv("REDIS_EXPIRATION", "1m")
	config.RedisExpiration, err = time.ParseDuration(redisExpiration)
	if err != nil {
		return config, fmt.Errorf("could not parse redis expiration: %w", err)
	}

	return config, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return strings.ToLower(value)
	}
	return fallback
}
