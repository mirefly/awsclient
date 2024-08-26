package awsclient

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/savaki/jq"
)

type secretsmanagerClient struct {
	Raw     *secretsmanager.Client
	timeout *timeout
}

const secretsManagerQueryTimeout = 5 * time.Second

func NewSecretsmanagerClient(cfg aws.Config) *secretsmanagerClient {
	return &secretsmanagerClient{
		Raw:     secretsmanager.NewFromConfig(cfg),
		timeout: newTimeout(10 * time.Second),
	}
}

func (sc *secretsmanagerClient) getSecretRawWithContext(ctx context.Context, secretName string) (string, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: &secretName,
	}

	result, err := sc.Raw.GetSecretValue(ctx, input)
	if err != nil {
		return "", err
	}

	secret := ""
	if result.SecretString != nil {
		secret = *result.SecretString
	}

	return secret, nil
}

func (sc *secretsmanagerClient) getSecretValueWithContext(ctx context.Context, secretName string, selector string) (string, error) {
	raw, err := sc.getSecretRawWithContext(ctx, secretName)
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(selector, ".") {
		selector = "." + selector
	}

	op, err := jq.Parse(selector)
	if err != nil {
		return "", err
	}
	bs, err := op.Apply([]byte(raw))
	return string(bs), err
}

func (sc *secretsmanagerClient) EasyGetSecretValue(secretName string, selector string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), secretsManagerQueryTimeout)
	defer cancel()

	return sc.getSecretValueWithContext(ctx, secretName, selector)
}
