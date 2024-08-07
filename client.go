package awsclient

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

type AWSClient struct {
	Config aws.Config

	dynamodbClient *dynamodbClient
	secretsmanager *secretsmanagerClient
	sesClient      *sesClient
}

func NewWithContext(cxt context.Context) (*AWSClient, error) {
	cfg, err := config.LoadDefaultConfig(cxt)
	if err != nil {
		return nil, err
	}

	return &AWSClient{Config: cfg}, nil
}

func New() (*AWSClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	return NewWithContext(ctx)
}

func (c *AWSClient) Dynamodb() *dynamodbClient {
	if c.dynamodbClient == nil {
		c.dynamodbClient = NewDynamodbClient(c.Config)
	}
	return c.dynamodbClient
}

func (c *AWSClient) Secretsmanager() *secretsmanagerClient {
	if c.secretsmanager == nil {
		c.secretsmanager = NewSecretsmanagerClient(c.Config)
	}
	return c.secretsmanager
}

func (c *AWSClient) Ses() *sesClient {
	if c.sesClient == nil {
		c.sesClient = NewSesClient(c.Config)
	}
	return c.sesClient
}
