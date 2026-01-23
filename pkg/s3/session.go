package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"

	"fmt"
)

func NewSession(cfg Config) (*session.Session, error) {
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("failed to validate s3 config: %w", err)
	}

	session, err := session.NewSession(&aws.Config{
		Endpoint:         aws.String(cfg.Endpoint),
		Region:           aws.String(cfg.Region),
		S3ForcePathStyle: aws.Bool(true),
		Credentials: credentials.NewStaticCredentials(
			cfg.AccessKeyID,
			cfg.AccessSecretKey,
			"",
		),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create s3 session: %w", err)
	}

	return session, nil
}

func NewExternalSession(cfg Config) (*session.Session, error) {
	if err := cfg.validateExternal(); err != nil {
		return nil, fmt.Errorf("failed to validate s3 config: %w", err)
	}

	session, err := session.NewSession(&aws.Config{
		Endpoint:         aws.String(cfg.ExternalEndpoint),
		Region:           aws.String(cfg.Region),
		S3ForcePathStyle: aws.Bool(cfg.ExternalS3ForcePathStyle),
		Credentials: credentials.NewStaticCredentials(
			cfg.AccessKeyID,
			cfg.AccessSecretKey,
			"",
		),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create s3 session: %w", err)
	}

	return session, nil
}
