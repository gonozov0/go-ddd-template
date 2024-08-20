package internal

import (
	"fmt"
	"os"
	"strings"
	"time"

	"go-echo-ddd-template/pkg/environment"
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
	Environment       environment.Type
	Port              string
	InterruptTimeout  time.Duration
	ReadHeaderTimeout time.Duration
	PprofPort         string
}

func loadServer() (Server, error) {
	var server Server

	server.Environment = environment.Type(getEnv("ENV_TYPE", string(environment.Testing)))
	server.Port = getEnv("SERVER_PORT", "8080")
	interruptTimeout, err := time.ParseDuration(getEnv("INTERRUPT_TIMEOUT", "2s"))
	if err != nil {
		return server, fmt.Errorf("could not parse interrupt timeout: %w", err)
	}
	server.InterruptTimeout = interruptTimeout
	readHeaderTimeout, err := time.ParseDuration(getEnv("READ_HEADER_TIMEOUT", "5s"))
	if err != nil {
		return server, fmt.Errorf("could not parse read header timeout: %w", err)
	}
	server.ReadHeaderTimeout = readHeaderTimeout
	server.PprofPort = getEnv("PPROF_PORT", "6060")

	return server, nil
}

type Sentry struct {
	DSN         string
	Environment environment.Type
}

func loadSentry() Sentry {
	var sentry Sentry

	sentry.Environment = environment.Type(getEnv("SENTRY_ENVIRONMENT", string(environment.Testing)))
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
