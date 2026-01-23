package sentry

import (
	"go-ddd-template/pkg/environment"
)

type Config struct {
	DSN         string
	Environment environment.Type
	Project     string
	Service     string
	Release     string
}
