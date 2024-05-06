package internal

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type Config struct {
	Server Server
	Sentry Sentry
	Redis  Redis
}

func LoadConfig() (Config, error) {
	var (
		config Config
		err    error
	)

	config.Server, err = loadServer()
	if err != nil {
		return config, fmt.Errorf("could not load server config: %w", err)
	}
	config.Sentry = loadSentry()
	config.Redis, err = loadRedis()
	if err != nil {
		return config, fmt.Errorf("could not load redis config: %w", err)
	}

	return config, nil
}

type Server struct {
	Port             string
	InterruptTimeout time.Duration
}

func loadServer() (Server, error) {
	var server Server

	server.Port = getEnv("SERVER_PORT", "8080")
	interruptTimeout, err := time.ParseDuration(getEnv("KILL_TIMEOUT", "2s"))
	if err != nil {
		return server, fmt.Errorf("could not parse kill timeout: %w", err)
	}
	server.InterruptTimeout = interruptTimeout

	return server, nil
}

type Sentry struct {
	DSN         string
	Environment string
}

func loadSentry() Sentry {
	var sentry Sentry

	sentry.Environment = getEnv("SENTRY_ENVIRONMENT", "development")
	sentry.DSN = getEnv("SENTRY_DSN", "")

	return sentry
}

type Redis struct {
	ClusterMode bool
	TLSEnabled  bool
	Address     string
	Username    string
	Password    string
	Expiration  time.Duration
}

func loadRedis() (Redis, error) {
	var redis Redis

	redis.ClusterMode = getEnv("REDIS_CLUSTER_MODE", "false") == "true"
	redis.TLSEnabled = getEnv("REDIS_TLS_ENABLED", "false") == "true"
	redis.Address = getEnv("REDIS_ADDRESS", "localhost:6379")
	redis.Username = getEnv("REDIS_USERNAME", "")
	redis.Password = getEnv("REDIS_PASSWORD", "")
	redisExpiration := getEnv("REDIS_EXPIRATION", "1m")
	expiration, err := time.ParseDuration(redisExpiration)
	if err != nil {
		return redis, fmt.Errorf("could not parse redis expiration: %w", err)
	}
	redis.Expiration = expiration

	return redis, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return strings.ToLower(value)
	}
	return fallback
}
