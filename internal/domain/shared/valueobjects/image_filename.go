package valueobjects

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	"errors"
)

type ImageFilename string

const (
	EmptyImageFilename   = ImageFilename("")
	ImageFilenamePattern = "%s_%s"
)

var (
	ErrInvalidImageFilename = errors.New("invalid image filename")
	allowedImageExtensions  = map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
		".bmp":  true,
		".svg":  true,
	}
)

func NewImageFilename(filename string) (ImageFilename, error) {
	if filename == "" {
		return EmptyImageFilename, nil
	}

	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		return "", fmt.Errorf("%w: filename must have an extension", ErrInvalidImageFilename)
	}

	if !allowedImageExtensions[ext] {
		return "", fmt.Errorf("%w: unsupported image extension: %s", ErrInvalidImageFilename, ext)
	}

	return ImageFilename(filename), nil
}

func NewRequiredImageFilenameForProduct(filename string, productID ProductID) (ImageFilename, error) {
	if filename == "" {
		return "", fmt.Errorf("%w: filename must not be empty", ErrInvalidImageFilename)
	}

	if productID.UUID() == uuid.Nil {
		return "", fmt.Errorf("%w: product id must not be empty", ErrInvalidProductID)
	}

	return NewImageFilename(fmt.Sprintf(ImageFilenamePattern, productID.String(), filename))
}

func (f ImageFilename) String() string {
	return string(f)
}

func (f ImageFilename) IsEmpty() bool {
	return f == EmptyImageFilename
}
