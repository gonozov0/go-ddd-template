package sqs

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	awssession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Session struct {
	SQS           *sqs.SQS
	QueueBaseName string
}

func NewSession(config Config) (Session, error) {
	creds := credentials.NewStaticCredentials(
		config.AccessKeyID,
		"unused",
		config.SessionToken,
	)

	awsConfig := &aws.Config{
		Region:      aws.String("yandex"),
		Endpoint:    aws.String("http://" + config.Endpoint),
		Credentials: creds,
	}

	awsSession, err := awssession.NewSession()
	if err != nil {
		return Session{}, fmt.Errorf("failed to init aws session: %w", err)
	}

	var session Session

	session.SQS = sqs.New(awsSession, awsConfig)
	session.QueueBaseName = fmt.Sprintf("http://%s/%s/%s", config.Endpoint, config.AccessKeyID, config.Environment)

	return session, nil
}
