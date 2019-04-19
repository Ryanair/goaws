package s3

import (
	"github.com/Ryanair/goaws/internal"

	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	ErrCodeSigningURL = "SigningURLErr"
)

type Error internal.Error

func (e Error) Error() string {
	return e.Message
}

func (e Error) Cause() error {
	return e.Causer
}

func wrapErr(err error, msg string) error {
	return Error(internal.WrapErr(err, msg))
}

func wrapErrWithCode(err error, msg, code string) error {
	return Error(internal.WrapErrWithCode(err, msg, code))
}

func (e Error) SigningFailed() bool {
	return internal.AnyEquals(e.Code, ErrCodeSigningURL)
}

func (e Error) BucketNotFound() bool {
	return internal.AnyEquals(e.Code, s3.ErrCodeNoSuchBucket)
}

func (e Error) KeyNotFound() bool {
	return internal.AnyEquals(e.Code, s3.ErrCodeNoSuchKey)
}

func (e Error) ResourceNotFound() bool {
	return internal.AnyEquals(e.Code,
		s3.ErrCodeNoSuchKey,
		s3.ErrCodeNoSuchBucket)
}

func (e Error) BucketAlreadyExists() bool {
	return internal.AnyEquals(e.Code,
		s3.ErrCodeBucketAlreadyExists)
}
