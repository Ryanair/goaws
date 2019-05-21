package s3

import (
	"io"
	"time"

	"github.com/Ryanair/goaws"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Client struct {
	s3 *s3.S3
}

func NewClient(cfg *goaws.Config, options ...func(*s3.S3)) *Client {
	cli := s3.New(cfg.Provider)
	for _, opt := range options {
		opt(cli)
	}

	return &Client{s3: cli}
}

func Endpoint(endpoint string) func(*s3.S3) {
	return func(s3 *s3.S3) {
		s3.Endpoint = endpoint
	}
}

func (c *Client) GeneratePutURL(bucket, key, contentType string, expire time.Duration) (string, error) {
	req, _ := c.s3.PutObjectRequest(&s3.PutObjectInput{
		Bucket:      &bucket,
		Key:         &key,
		ContentType: &contentType,
	})

	url, err := req.Presign(expire)
	if err != nil {
		return "", wrapErrWithCode(err, ErrCodeSigningURL, "signing put url failed")
	}

	return url, nil
}

func (c *Client) GeneratePutURLWithMetadata(bucket, key, contentType string, expire time.Duration, metadata map[string]*string) (string, error) {
	req, _ := c.s3.PutObjectRequest(&s3.PutObjectInput{
		Bucket:      &bucket,
		Key:         &key,
		ContentType: &contentType,
		Metadata:    metadata,
	})

	url, err := req.Presign(expire)
	if err != nil {
		return "", wrapErrWithCode(err, ErrCodeSigningURL, "signing put url with metadata failed")
	}

	return url, nil
}

func (c *Client) DeleteObject(bucket, key string) error {
	if _, err := c.s3.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}); err != nil {
		return wrapErr(err, "delete object failed")
	}

	return nil
}

func (c *Client) GetObject(bucket, key string) (io.ReadCloser, error) {
	out, err := c.s3.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, wrapErr(err, "get object failed")
	}

	return out.Body, nil
}

func (c *Client) PutObject(bucket, key string, body io.ReadSeeker) error {
	_, err := c.s3.PutObject(&s3.PutObjectInput{
		Body:   body,
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return wrapErr(err, "put object failed")
	}

	return nil
}

func (c *Client) PutObjectWithMetadata(bucket, key string, body io.ReadSeeker, metadata map[string]*string) error {
	_, err := c.s3.PutObject(&s3.PutObjectInput{
		Body:     body,
		Bucket:   aws.String(bucket),
		Key:      aws.String(key),
		Metadata: metadata,
	})
	if err != nil {
		return wrapErr(err, "put object with metadata failed")
	}

	return nil
}

func (c *Client) GetObjectMetadata(bucket, key string) (map[string]*string, error) {
	out, err := c.s3.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, wrapErr(err, "get object metadata failed")
	}

	return out.Metadata, nil
}
