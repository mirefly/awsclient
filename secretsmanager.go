package awsclient

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
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

func (sc *secretsmanagerClient) getSecretJSONWithContext(ctx context.Context, secretName string) (map[string]string, error) {
	v := make(map[string]string)

	secret, err := sc.getSecretRawWithContext(ctx, secretName)
	if err != nil {
		return v, err
	}

	err = json.Unmarshal([]byte(secret), &v)
	return v, err
}

func (sc *secretsmanagerClient) getSecretValueWithContext(ctx context.Context, secretName string, key string) (string, error) {
	secretJSON, err := sc.getSecretJSONWithContext(ctx, secretName)
	if err != nil {
		return "", err
	}

	value, ok := secretJSON[key]
	if !ok {
		return "", fmt.Errorf("value for key '%s' doesn't exist", key)
	}

	return value, nil
}

func (sc *secretsmanagerClient) EasyGetSecretValue(secretName string, key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), secretsManagerQueryTimeout)
	defer cancel()

	return sc.getSecretValueWithContext(ctx, secretName, key)
}
