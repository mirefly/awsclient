package awsclient

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type dynamodbClient struct {
	Raw     *dynamodb.Client
	timeout *timeout
}

func NewDynamodbClient(cfg aws.Config) *dynamodbClient {
	return &dynamodbClient{
		Raw:     dynamodb.NewFromConfig(cfg),
		timeout: newTimeout(10 * time.Second),
	}
}

func (dc *dynamodbClient) Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	cxt, cancel := context.WithTimeout(context.Background(), dc.timeout.Value())
	defer cancel()
	return dc.Raw.Query(cxt, input)
}

func (dc *dynamodbClient) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	cxt, cancel := context.WithTimeout(context.Background(), dc.timeout.Value())
	defer cancel()
	return dc.Raw.PutItem(cxt, input)
}
