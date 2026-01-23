package internal

import (
	"fmt"
	"time"

	"go-ddd-template/pkg/auth"

	"go-ddd-template/pkg/db/postgres"
	"go-ddd-template/pkg/db/redis"
	"go-ddd-template/pkg/db/ydb"
	"go-ddd-template/pkg/environment"
	"go-ddd-template/pkg/envutils"
	imagestorage "go-ddd-template/pkg/image_storage"
	"go-ddd-template/pkg/logger"
	"go-ddd-template/pkg/logger/sentry"
	"go-ddd-template/pkg/s3"
	"go-ddd-template/pkg/sqs"
	"go-ddd-template/pkg/traces"
)

type Config struct {
	Server       Server
	Crons        Crons
	Consumers    Consumers
	Auth         auth.Config
	Postgres     postgres.Config
	YDB          ydb.Config
	Logger       logger.Config
	Traces       traces.Config
	Redis        redis.Config
	SQS          sqs.Config
	ImageStorage imagestorage.Config
}

type Consumers struct {
	HTTPPort string
}

func loadConsumers() Consumers {
	var consumers Consumers

	consumers.HTTPPort = envutils.GetEnv("CONSUMER_HTTP_PORT", "8083")

	return consumers
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

	config.Crons, err = loadCrons()
	if err != nil {
		return config, fmt.Errorf("could not load crons config: %w", err)
	}

	config.Auth, err = auth.LoadConfig()
	if err != nil {
		return config, fmt.Errorf("could not load auth config: %w", err)
	}

	config.Postgres = loadPostgres()
	config.YDB = loadYDB()
	config.Logger = loadLogger(config.Server.Environment)
	config.Consumers = loadConsumers()
	config.SQS = loadSQS(config.Server.Environment)

	config.Traces, err = traces.LoadConfig()
	if err != nil {
		return config, fmt.Errorf("could not load traces config: %w", err)
	}

	config.Redis = loadRedis()

	config.ImageStorage = loadImageStorage()

	return config, nil
}

type Server struct {
	Environment       environment.Type
	HTTPPort          string
	GRPCPort          string
	InterruptTimeout  time.Duration
	ReadHeaderTimeout time.Duration
	PprofPort         string
	MetricPort        string
}

func loadServer() (Server, error) {
	var (
		server Server
		err    error
	)

	server.Environment, err = environment.GetType(envutils.GetEnv("ENV_TYPE", string(environment.Local)))
	if err != nil {
		return server, fmt.Errorf("could not get environment type: %w", err)
	}

	server.HTTPPort = envutils.GetEnv("HTTP_PORT", "8080")
	server.GRPCPort = envutils.GetEnv("GRPC_PORT", "8081")

	interruptTimeout, err := time.ParseDuration(envutils.GetEnv("INTERRUPT_TIMEOUT", "2s"))
	if err != nil {
		return server, fmt.Errorf("could not parse interrupt timeout: %w", err)
	}

	server.InterruptTimeout = interruptTimeout

	readHeaderTimeout, err := time.ParseDuration(envutils.GetEnv("READ_HEADER_TIMEOUT", "5s"))
	if err != nil {
		return server, fmt.Errorf("could not parse read header timeout: %w", err)
	}

	server.ReadHeaderTimeout = readHeaderTimeout
	server.PprofPort = envutils.GetEnv("PPROF_PORT", "6060")
	// 3400-3499 is open ports for monitoring  https://m.yandex-team.ru/docs/network_access#access-pull
	server.MetricPort = envutils.GetEnv("METRIC_PORT", "")

	return server, nil
}

type Crons struct {
	HandleCraetedOrdersInterval time.Duration
	HTTPPort                    string
}

func loadCrons() (Crons, error) {
	var (
		crons Crons
		err   error
	)

	crons.HandleCraetedOrdersInterval, err = time.ParseDuration(
		envutils.GetEnv("HANDLE_CREATED_ORDERS_INTERVAL", "30s"),
	)
	if err != nil {
		return crons, fmt.Errorf("could not parse handle created orders interval: %w", err)
	}

	crons.HTTPPort = envutils.GetEnv("CRON_HTTP_PORT", "8082")

	return crons, nil
}

func loadPostgres() postgres.Config {
	return postgres.Config{
		Hosts:       envutils.SplitEnv("POSTGRES_HOSTS", "localhost"),
		Port:        envutils.GetEnv("POSTGRES_PORT", "5432"),
		User:        envutils.GetEnv("POSTGRES_USER", "postgres"),
		Password:    envutils.GetEnv("POSTGRES_PASSWORD", "postgres"),
		Database:    envutils.GetEnv("POSTGRES_DATABASE", "postgres"),
		SSL:         envutils.GetEnvBool("POSTGRES_SSL", false),
		SSLRootCert: envutils.GetEnv("POSTGRES_SSL_ROOT_CERT", ""),
	}
}

func loadYDB() ydb.Config {
	return ydb.Config{
		Endpoint: envutils.GetEnv("YDB_ENDPOINT", "localhost:2135"),
		Database: envutils.GetEnv("YDB_DATABASE", "local"),
		Token:    envutils.GetEnv("YDB_TOKEN", ""),
	}
}

func loadLogger(envType environment.Type) logger.Config {
	return logger.Config{
		EnableJSONFormat: envutils.GetEnv("DEPLOY_UNIT_ID", "") != "",
		LogLevel:         logger.GetLogLevel(envutils.GetEnv("LOG_LEVEL", "info")),
		Sentry:           loadSentry(envType),
	}
}

func loadSentry(envType environment.Type) sentry.Config {
	return sentry.Config{
		DSN:         envutils.GetEnv("SENTRY_DSN", ""),
		Environment: envType,
		Project:     envutils.GetEnv("PROJECT_ID", ""),
		Service:     envutils.GetEnv("UNIT_ID", ""),
		Release:     envutils.GetEnv("RELEASE_VERSION", ""),
	}
}

func loadRedis() redis.Config {
	return redis.Config{
		TlsEnabled:  envutils.GetEnvBool("REDIS_TLS_ENABLED", false),
		Addrs:       envutils.SplitEnv("REDIS_ADDRS", "localhost:6379"),
		Username:    envutils.GetEnv("REDIS_USERNAME", ""),
		Password:    envutils.GetEnv("REDIS_PASSWORD", ""),
		TlsRootCert: envutils.GetEnv("REDIS_TLS_ROOT_CERT", ""),
	}
}

func loadSQS(envType environment.Type) sqs.Config {
	defaultEndpoint := fmt.Sprintf("localhost:%s", envutils.GetEnv("SQS_PORT", "9324"))

	return sqs.Config{
		Endpoint:     envutils.GetEnv("SQS_ENDPOINT", defaultEndpoint),
		AccessKeyID:  envutils.GetEnv("SQS_ACCESS_KEY_ID", "000011112222"),
		SessionToken: envutils.GetEnv("SQS_SESSION_TOKEN", ""),
		Environment:  envType,
	}
}

func loadImageStorage() imagestorage.Config {
	return imagestorage.Config{
		ImageBucketName: envutils.GetEnv("S3_IMAGE_BUCKET_NAME", "test-bucket"),
		S3Config:        s3.LoadConfig(),
	}
}
