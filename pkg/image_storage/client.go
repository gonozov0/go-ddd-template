package imagestorage

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	awss3 "github.com/aws/aws-sdk-go/service/s3"

	"go-ddd-template/internal/domain/shared/valueobjects"
	"go-ddd-template/pkg/s3"
)

type Client struct {
	s3              *awss3.S3
	imageBucketName string
	endpoint        string
}

func NewClient(ctx context.Context, cfg Config) (*Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid image storage config: %w", err)
	}

	session, err := s3.NewSession(cfg.S3Config)
	if err != nil {
		return nil, fmt.Errorf("failed to create s3 session: %w", err)
	}

	c := &Client{
		s3:              awss3.New(session),
		imageBucketName: cfg.ImageBucketName,
		endpoint:        cfg.S3Config.Endpoint,
	}

	if err := c.createImageBucketIfNotExists(ctx); err != nil {
		return nil, fmt.Errorf("failed to create image bucket if not exists: %w", err)
	}

	return c, nil
}

func (c *Client) GetImageUploadURL(filename valueobjects.ImageFilename) (valueobjects.ImageURL, error) {
	req, _ := c.s3.PutObjectRequest(&awss3.PutObjectInput{
		Bucket: aws.String(c.imageBucketName),
		Key:    aws.String(filename.String()),
		ACL:    aws.String(ImageBucketACL),
	})

	url, err := req.Presign(defaultPresignTTL)
	if err != nil {
		return "", fmt.Errorf("failed to get upload image url: %w", err)
	}

	return valueobjects.NewImageURL(url), nil
}

func (c *Client) ConfirmImageUpload(
	ctx context.Context,
	filename valueobjects.ImageFilename,
) (valueobjects.ImageURL, error) {
	if err := c.checkImageExists(ctx, filename); err != nil {
		return valueobjects.EmptyImageURL, fmt.Errorf("failed to check image exists: %w", err)
	}

	return valueobjects.NewImageURL(fmt.Sprintf("%s/%s/%s", c.endpoint, c.imageBucketName, filename.String())), nil
}

func (c *Client) DeleteImage(ctx context.Context, filename valueobjects.ImageFilename) error {
	_, err := c.s3.DeleteObjectWithContext(ctx, &awss3.DeleteObjectInput{
		Bucket: aws.String(c.imageBucketName),
		Key:    aws.String(filename.String()),
	})
	if err != nil {
		return fmt.Errorf("failed to delete image: %w", err)
	}

	return nil
}

func (c *Client) checkImageExists(ctx context.Context, filename valueobjects.ImageFilename) error {
	headObjectInput := &awss3.HeadObjectInput{
		Bucket: aws.String(c.imageBucketName),
		Key:    aws.String(filename.String()),
	}

	_, err := c.s3.HeadObjectWithContext(ctx, headObjectInput)
	if err != nil {
		return fmt.Errorf(
			"file %s does not exist in bucket %s: %w",
			filename.String(),
			c.imageBucketName,
			ErrImageNotFound,
		)
	}

	return nil
}

func (c *Client) createImageBucketIfNotExists(ctx context.Context) error {
	listBucketsOutput, err := c.s3.ListBucketsWithContext(ctx, &awss3.ListBucketsInput{})
	if err != nil {
		return fmt.Errorf("failed to list buckets: %w", err)
	}

	bucketExists := false

	for _, bucket := range listBucketsOutput.Buckets {
		if bucket.Name != nil && *bucket.Name == c.imageBucketName {
			bucketExists = true
			break
		}
	}

	if !bucketExists {
		if _, err := c.s3.CreateBucketWithContext(ctx, &awss3.CreateBucketInput{
			Bucket: aws.String(c.imageBucketName),
			ACL:    aws.String(ImageBucketACL),
		}); err != nil {
			return fmt.Errorf("failed to create bucket %s: %w", c.imageBucketName, err)
		}
	}

	return nil
}
