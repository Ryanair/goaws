package s3

import (
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
	"github.com/ryanair/goaws"
)

const (
	defaultExpiration = 30 * time.Minute
)

type Client struct {
	*s3.S3
}

func NewClient(cfg *goaws.Config) *Client {
	return &Client{S3: s3.New(cfg.Provider)}
}

func (c *Client) GeneratePreSignedPutURL(bucket, key, contentType string) (string, error) {
	req, _ := c.PutObjectRequest(&s3.PutObjectInput{
		Bucket:      &bucket,
		Key:         &key,
		ContentType: &contentType,
	})

	url, err := req.Presign(defaultExpiration)
	if err != nil {
		return "", newError(errors.Wrap(err, "signing url failed").Error(), SigningURLErrCode)
	}

	return url, nil
}
