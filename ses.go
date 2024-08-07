package awsclient

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

type sesClient struct {
	Raw     *ses.Client
	timeout *timeout
}

func NewSesClient(cfg aws.Config) *sesClient {
	return &sesClient{
		Raw:     ses.NewFromConfig(cfg),
		timeout: newTimeout(10 * time.Second),
	}
}

func (sc *sesClient) SendEmail(input *ses.SendEmailInput) (*ses.SendEmailOutput, error) {
	cxt, cancel := context.WithTimeout(context.Background(), sc.timeout.Value())
	defer cancel()
	return sc.Raw.SendEmail(cxt, input)
}

func (sc *sesClient) SendSimpleEmail(from, to, subject, body string) (*ses.SendEmailOutput, error) {
	charset := "UTF-8"
	body = strings.Join(strings.Split(body, "\n"), "<br>")

	input := &ses.SendEmailInput{
		Source: &from,
		Destination: &types.Destination{
			ToAddresses: []string{to},
		},
		Message: &types.Message{
			Subject: &types.Content{
				Data: &subject,
			},
			Body: &types.Body{
				Html: &types.Content{
					Data:    &body,
					Charset: &charset,
				},
			},
		},
	}
	return sc.SendEmail(input)
}
