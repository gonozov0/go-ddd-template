package valueobjects

import (
	"errors"
	"fmt"
)

type ImageURL string

var ErrInvalidImageURL = errors.New("invalid image url")

const EmptyImageURL = ImageURL("")

func NewImageURL(url string) ImageURL {
	if url == "" {
		return EmptyImageURL
	}

	return ImageURL(url)
}

func NewRequiredImageURL(url string) (ImageURL, error) {
	if url == "" {
		return "", fmt.Errorf("%w: url must not be empty", ErrInvalidImageURL)
	}

	return NewImageURL(url), nil
}

func (url ImageURL) String() string {
	return string(url)
}

func (url ImageURL) IsEmpty() bool {
	return url == EmptyImageURL
}
