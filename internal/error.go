package internal

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/pkg/errors"
)

const (
	ErrCodeUnknownBehaviour = "ErrUnknownBehaviour"
)

type Error struct {
	Message string
	Code    string
	Causer  error
}

func NewError(message, code string, causer error) Error {
	return Error{
		Message: message,
		Code:    code,
		Causer:  causer,
	}
}

// WrapErrWithCode should be used to apply behavior to wrapped error, the behavior is specified by providing error code
func WrapErrWithCode(err error, msg, code string) Error {
	wrappedErrMsg := errors.Wrap(err, msg).Error()
	return NewError(wrappedErrMsg, code, err)
}

// WrapErr should be used to wrap AWS SDK errors or errors with unknown behavior
func WrapErr(err error, msg string) Error {
	if aerr, ok := err.(awserr.Error); ok {
		return WrapErrWithCode(err, msg, aerr.Code())
	}

	return WrapErrWithCode(err, msg, ErrCodeUnknownBehaviour)
}

func AnyEquals(origCode string, errCodes ...string) bool {
	for _, ec := range errCodes {
		if origCode == ec {
			return true
		}
	}

	return false
}
