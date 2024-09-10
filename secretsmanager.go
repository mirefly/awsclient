package awsclient

import (
	"context"
	"encoding/json"
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

func (sc *secretsmanagerClient) getSecretRawWithContext(ctx context.Context, secretName string) ([]byte, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: &secretName,
	}

	result, err := sc.Raw.GetSecretValue(ctx, input)
	if err != nil {
		return nil, err
	}

	secret := ""
	if result.SecretString != nil {
		secret = *result.SecretString
	}

	return []byte(secret), nil
}

// EasyGetSecretValue is a helper function to get secret value as byte slice
//
// If selector is empty, it will return the raw secret value as byte slice.
//
// If selector is not empty, it will expect the secret value to be a JSON string and will return the selected part as byte slice.
// If the selected field is a string, double quotes will be kept. `EasyGetSecretValueS` should be used instead.
func (sc *secretsmanagerClient) EasyGetSecretValue(secretName string, selector string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), secretsManagerQueryTimeout)
	defer cancel()

	raw, err := sc.getSecretRawWithContext(ctx, secretName)
	if err != nil {
		return nil, err
	}

	if selector == "" {
		return raw, nil
	}

	if !strings.HasPrefix(selector, ".") {
		selector = "." + selector
	}

	op, err := jq.Parse(selector)
	if err != nil {
		return nil, err
	}
	return op.Apply([]byte(raw))
}

// EasyGetSecretValueS is a helper function to get secret value as string
//
// If selector is empty, it will return the secret value as string.
//
// If selector is not empty, it will expect the secret value to be a JSON string and will return the selected value as string.
// If the selected field is not string, it will return error.
func (sc *secretsmanagerClient) EasyGetSecretValueS(secretName string, selector string) (string, error) {
	bs, err := sc.EasyGetSecretValue(secretName, selector)
	if err != nil {
		return "", err
	}

	if selector == "" {
		return string(bs), nil
	}

	var s string
	if err := json.Unmarshal(bs, &s); err != nil {
		return "", err
	}

	return s, nil
}
