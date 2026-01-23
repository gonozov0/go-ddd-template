package sqs

import "go-ddd-template/pkg/environment"

type Config struct {
	Endpoint     string
	AccessKeyID  string
	SessionToken string
	Environment  environment.Type
}
