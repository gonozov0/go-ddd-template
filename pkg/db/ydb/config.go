package ydb

import "errors"

type Config struct {
	Endpoint string
	Database string
	Token    string
}

func (c Config) validate() error {
	if c.Endpoint == "" {
		return errors.New("Endpoint is required")
	}

	return nil
}

func (c Config) DSN() string {
	return "grpc://" + c.Endpoint + "/" + c.Database
}
