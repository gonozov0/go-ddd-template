package environment

import (
	"strings"

	"errors"
)

type Type string

const (
	Production Type = "production"
	Testing    Type = "testing"
	Dev        Type = "dev"
	Local      Type = "local"
)

var ErrInvalidType = errors.New("invalid environment type")

func GetType(envType string) (Type, error) {
	switch strings.ToLower(envType) {
	case string(Production):
		return Production, nil
	case string(Testing):
		return Testing, nil
	case string(Dev):
		return Dev, nil
	case string(Local):
		return Local, nil
	}

	return "", ErrInvalidType
}

func (t Type) IsLocal() bool {
	return t == Local
}
