package s3

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/ryanair/goaws/internal"
)

const (
	SigningURLErrCode = "SigningURLErr"
)

type Error internal.Error

func (e Error) Error() string {
	return e.Message
}

func newErr(err error, code, msg string) error {
	return Error(internal.WrapErr(err, code, msg))
}

func newOpsErr(err error, msg string) error {
	return Error(internal.WrapOpsErr(err, msg))
}

func (e Error) SigningFailed() bool {
	return internal.AnyEquals(e.Code, SigningURLErrCode)
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
