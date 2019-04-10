package goaws

import (
	"os"

	"github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/pkg/errors"
)

type Config struct {
	Provider client.ConfigProvider
}

func NewConfig(options ...func(*aws.Config)) (*Config, error) {
	config, err := newConfig(options...)
	if err != nil {
		return &Config{}, errors.Wrap(err, "unable to create config")
	}
	return config, nil
}

func newConfig(options ...func(*aws.Config)) (*Config, error) {
	awsConfig := &aws.Config{}
	for _, opt := range options {
		opt(awsConfig)
	}

	if awsConfig.Region == nil || *awsConfig.Region == "" {
		regionEnvVar := os.Getenv("AWS_REGION")
		if regionEnvVar == "" {
			return &Config{}, errors.New("AWS_REGION environment variable not found")
		}
		Region(regionEnvVar)(awsConfig)
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return &Config{}, errors.Wrap(err, "unable to create AWS session")
	}

	return &Config{Provider: sess}, nil
}

func Region(region string) func(*aws.Config) {
	return func(c *aws.Config) {
		c.Region = aws.String(region)
	}
}

func MaxRetries(max int) func(*aws.Config) {
	return func(c *aws.Config) {
		c.MaxRetries = aws.Int(max)
	}
}

func Credentials(id, secret, token string) func(*aws.Config) {
	return func(c *aws.Config) {
		c.WithCredentials(credentials.NewStaticCredentials(id, secret, token))
	}
}
