package imagestorage

import (
	"time"

	"errors"

	"go-ddd-template/pkg/s3"
)

const (
	defaultPresignTTL = 5 * time.Minute
	ImageBucketACL    = "public-read"
)

type Config struct {
	ImageBucketName string
	S3Config        s3.Config
}

func (cfg Config) Validate() error {
	if cfg.ImageBucketName == "" {
		return errors.New("image bucket name is required")
	}

	return nil
}
