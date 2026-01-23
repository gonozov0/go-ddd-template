package s3

import (
	"fmt"
	"os"

	"go-ddd-template/pkg/envutils"
)

type Config struct {
	Endpoint                 string
	ExternalEndpoint         string
	ExternalS3ForcePathStyle bool
	AccessKeyID              string
	AccessSecretKey          string
	Region                   string
}

func (c *Config) validate() error {
	if c.Endpoint == "" {
		return fmt.Errorf("s3 endpoint is required")
	}

	if c.AccessKeyID == "" {
		return fmt.Errorf("s3 access key id is required")
	}

	if c.AccessSecretKey == "" {
		return fmt.Errorf("s3 secret access key is required")
	}

	if c.Region == "" {
		return fmt.Errorf("s3 region is required")
	}

	return nil
}

func (c *Config) validateExternal() error {
	if c.ExternalEndpoint == "" {
		return fmt.Errorf("s3 external endpoint is required")
	}

	return c.validate()
}

func LoadConfig() Config {
	// Yandex Cloud S3 (Default)
	if os.Getenv("S3_ENDPOINT") != "" {
		return Config{
			Endpoint:         envutils.GetEnv("S3_ENDPOINT", ""),
			ExternalEndpoint: envutils.GetEnv("S3_EXTERNAL_ENDPOINT", ""),
			// External urls for YC S3 only available with virtual-hosted style: https://{{bucket}}.s3-private.mds.yandex.net/{{key}}
			ExternalS3ForcePathStyle: envutils.GetEnvBool("S3_EXTERNAL_S3_FORCE_PATH_STYLE", false),
			AccessKeyID:              envutils.GetEnv("S3_ACCESS_KEY_ID", ""),
			AccessSecretKey:          envutils.GetEnv("S3_ACCESS_SECRET_KEY", ""),
			Region:                   envutils.GetEnv("S3_REGION", "ru-central1"),
		}
	}

	// S3MDS Recipe (Integration Tests)
	if os.Getenv("S3MDS_PORT") != "" {
		endpoint := envutils.GetEnv(
			"S3_ENDPOINT",
			fmt.Sprintf("http://localhost:%s", os.Getenv("S3MDS_PORT")),
		)

		return Config{
			Endpoint:                 endpoint,
			ExternalEndpoint:         endpoint,
			ExternalS3ForcePathStyle: true,
			AccessKeyID:              envutils.GetEnv("S3_ACCESS_KEY_ID", "test"),
			AccessSecretKey:          envutils.GetEnv("S3_ACCESS_SECRET_KEY", "test"),
			Region:                   envutils.GetEnv("S3_REGION", "eu-central-1"),
		}
	}

	// Localstack in Docker (Debugging, Integration Tests)
	return Config{
		Endpoint:                 envutils.GetEnv("S3_ENDPOINT", "http://localhost:4566"),
		ExternalEndpoint:         envutils.GetEnv("S3_EXTERNAL_ENDPOINT", "http://localhost:4566"),
		ExternalS3ForcePathStyle: envutils.GetEnvBool("S3_EXTERNAL_S3_FORCE_PATH_STYLE", true),
		AccessKeyID:              envutils.GetEnv("S3_ACCESS_KEY_ID", "test"),
		AccessSecretKey:          envutils.GetEnv("S3_ACCESS_SECRET_KEY", "test"),
		Region:                   envutils.GetEnv("S3_REGION", "eu-central-1"),
	}
}
