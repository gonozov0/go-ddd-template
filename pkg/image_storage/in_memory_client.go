package imagestorage

import (
	"context"

	"github.com/brianvoe/gofakeit/v6"

	"fmt"

	"go-ddd-template/internal/domain/shared/valueobjects"
)

type InMemoryClient struct {
	files map[valueobjects.ImageFilename]bool
}

func NewInMemoryClient(ctx context.Context) *InMemoryClient {
	return &InMemoryClient{
		files: make(map[valueobjects.ImageFilename]bool),
	}
}

func (c *InMemoryClient) GetImageUploadURL(filename valueobjects.ImageFilename) (valueobjects.ImageURL, error) {
	c.files[filename] = true

	return valueobjects.NewImageURL(gofakeit.URL()), nil
}

func (c *InMemoryClient) ConfirmImageUpload(
	ctx context.Context,
	filename valueobjects.ImageFilename,
) (valueobjects.ImageURL, error) {
	if !c.files[filename] {
		return valueobjects.EmptyImageURL, fmt.Errorf("file %s: %w", filename.String(), ErrImageNotFound)
	}

	return valueobjects.NewImageURL(gofakeit.URL()), nil
}

func (c *InMemoryClient) DeleteImage(ctx context.Context, filename valueobjects.ImageFilename) error {
	delete(c.files, filename)

	return nil
}
